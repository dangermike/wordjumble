package main

import (
	"bufio"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/DataDog/zstd"

	"github.com/dangermike/wordjumble/arraytrie"
	"github.com/dangermike/wordjumble/logging"
	"github.com/dangermike/wordjumble/maptrie"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	//go:embed dicts/*
	f       embed.FS
	newline = []byte("\n")
)

func main() {
	app := &cli.App{
		Name:   "wordjumble",
		Usage:  "permute letters against a dictionary",
		Action: appMain,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "use-array",
				Value: false,
				Usage: "Use the arrayTrie implementation instead of the mapTrie",
			},
			&cli.StringFlag{
				Name:    "dict",
				Aliases: []string{"d"},
				Value:   "2of12inf",
				Usage:   "Name of the dictionary to use",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Get wordy with those words",
			},
			&cli.BoolFlag{
				Name:    "consume",
				Aliases: []string{"c"},
				Usage:   "Consume letters (only use each letter once)",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Use all letters",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "Show available dictionaries",
				Action: listMain,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func listMain(c *cli.Context) error {
	entries, err := f.ReadDir("dicts")
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		fmt.Println(strings.TrimSuffix(entry.Name(), ".zst"))
	}
	return nil
}

func appMain(c *cli.Context) error {
	var t trie
	logger := zap.NewNop()
	if c.Bool("verbose") {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.DisableCaller = true
		logger, _ = cfg.Build()
	}

	ctx := logging.NewContext(c.Context, logger)

	if c.Bool("use-array") {
		t = arraytrie.New()
		logger.Debug("using array trie")
	} else {
		t = maptrie.New()
		logger.Debug("using map trie")
	}

	if err := loadWords(ctx, c.String("dict"), t.Load); err != nil {
		return err
	}

	args := c.Args()
	if 0 < args.Len() {
		return runWords(ctx, args.Slice(), t, c.Bool("consume"), c.Bool("all"))
	}
	message.NewPrinter(language.English).Printf("Loaded dictionary %s: %d words\n", c.String("dict"), t.Count())
	return runREPL(ctx, t, c.Bool("consume"), c.Bool("all"))
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

	for _, outword := range t.PermuteAll([]byte(word), consume) {
		if all && len(word) != len(outword) {
			continue
		}
		os.Stdout.Write(outword)
		os.Stdout.Write(newline)
		cnt++
	}
	logging.FromContext(ctx).Debug("permuted", zap.Int("count", cnt), zap.Duration("duration", time.Since(start)))
}

func loadWords(ctx context.Context, dict string, callback func(word []byte) bool) error {
	f, err := f.Open(fmt.Sprintf("dicts/%s.zst", dict))
	if err != nil {
		return err
	}
	reader := zstd.NewReader(f)
	defer reader.Close()
	breader := bufio.NewReader(reader)
	cnt := 0
	start := time.Now()
	for {
		line, isprefix, err := breader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if isprefix {
			return errors.New("Buffer not long enough to hold line")
		}
		if !callback(line) {
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
