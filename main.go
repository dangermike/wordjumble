package main

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"

	"github.com/dangermike/wordjumble/arraytrie"
	"github.com/dangermike/wordjumble/logging"
	"github.com/dangermike/wordjumble/maptrie"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	//go:embed dicts/*
	f embed.FS
)

func main() {
	cmdPermute := &cobra.Command{
		Use:  "permute permute letters against a dictionary",
		RunE: appMain,
	}
	AddFlags(cmdPermute.Flags())

	cmdList := &cobra.Command{
		Use:  "list Show available dictionaries",
		RunE: listMain,
	}

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(cmdPermute, cmdList)

	// use the default cmd if no cmd is given
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if (err == nil || strings.HasPrefix(err.Error(), "unknown command ")) && cmd == rootCmd && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{cmdPermute.Name()}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func listMain(cmd *cobra.Command, args []string) error {
	entries, err := f.ReadDir("dicts")
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		fmt.Println(strings.TrimSuffix(entry.Name(), ".zst"))
	}
	return nil
}

func appMain(cmd *cobra.Command, args []string) error {
	var t trie
	logger := zap.NewNop()

	cfg, err := GetConfig(cmd.Flags())
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true

	if cfg.Verbose {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.DisableCaller = true
		logger, _ = cfg.Build()
	}

	ctx := logging.NewContext(cmd.Context(), logger)

	if cfg.UseArray {
		t = arraytrie.New()
		logger.Debug("using array trie")
	} else {
		t = maptrie.New()
		logger.Debug("using map trie")
	}

	if err := loadWords(ctx, cfg.Dict, t.Load); err != nil {
		return err
	}

	if 0 < len(args) {
		return runWords(ctx, args, t, cfg.Consume, cfg.All)
	}
	message.NewPrinter(language.English).Printf("Loaded dictionary %s: %d words\n", cfg.Dict, t.Count())
	return runREPL(ctx, t, cfg.Consume, cfg.All)
}

func runREPL(ctx context.Context, t trie, consume bool, all bool) error {
	s := bufio.NewScanner(os.Stdin)
	os.Stdout.WriteString("words> ")
	for s.Scan() {
		if s.Text() == "" {
			return nil
		}
		runWord(ctx, s.Text(), t, consume, all)
		os.Stdout.WriteString("words> ")
	}
	if s.Err() == io.EOF {
		return nil
	}
	return s.Err()
}

func runWords(ctx context.Context, words []string, t trie, consume bool, all bool) error {
	for i, word := range words {
		if i > 0 {
			fmt.Println("------")
		}
		runWord(ctx, word, t, consume, all)
	}
	return nil
}

func runWord(ctx context.Context, word string, t trie, consume bool, all bool) {
	start := time.Now()
	cnt := 0

	bout := bufio.NewWriter(os.Stdout)
	for _, outword := range t.PermuteAll([]byte(word), consume) {
		if all && len(word) != len(outword) {
			continue
		}
		bout.Write(outword)
		bout.WriteByte('\n')
		cnt++
	}
	bout.Flush()
	logging.FromContext(ctx).Debug("permuted", zap.String("word", word), zap.Int("count", cnt), zap.Duration("duration", time.Since(start)))
}

func loadWords(ctx context.Context, dict string, callback func(word []byte) bool) error {
	f, err := f.Open(fmt.Sprintf("dicts/%s.zst", dict))
	if err != nil {
		return err
	}
	reader, err := zstd.NewReader(f)
	if err != nil {
		return err
	}
	defer reader.Close()
	scn := bufio.NewScanner(reader)
	cnt := 0
	start := time.Now()
	for scn.Scan() {
		if !callback(scn.Bytes()) {
			break
		}
		cnt++
	}

	logging.FromContext(ctx).Debug(
		"loaded",
		zap.Int("count", cnt),
		zap.Duration("duration", time.Since(start)),
		zap.String("dictionary", dict),
	)
	return nil
}

type trie interface {
	Load(word []byte) bool
	LoadString(word string) bool
	Count() int
	PermuteAll(letters []byte, consume bool) [][]byte
	Contains(letters []byte) bool
	ContainsString(word string) bool
}
