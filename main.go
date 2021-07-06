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
		t = &atWrapper{}
		logger.Debug("using array trie")
	} else {
		t = &mtWrapper{}
		logger.Debug("using map trie")
	}

	if err := loadWords(ctx, c.String("dict"), t.Load); err != nil {
		return err
	}

	args := c.Args()
	if 0 < args.Len() {
		return runWords(ctx, args.Slice(), t)
	}
	message.NewPrinter(language.English).Printf("Loaded dictionary %s: %d words\n", c.String("dict"), t.Count())
	return runREPL(ctx, t)
}

func runREPL(ctx context.Context, t trie) error {
	s := bufio.NewScanner(os.Stdin)
	os.Stdout.WriteString("words> ")
	for s.Scan() {
		runWord(ctx, s.Text(), t)
		os.Stdout.WriteString("words> ")
	}
	if s.Err() == io.EOF {
		return nil
	}
	return s.Err()
}

func runWords(ctx context.Context, words []string, t trie) error {
	for i, word := range words {
		if i > 0 {
			fmt.Println("------")
		}
		runWord(ctx, word, t)
	}
	return nil
}

func runWord(ctx context.Context, word string, t trie) {
	start := time.Now()
	cnt := 0

	for _, word := range t.PermuteAll([]byte(word)) {
		os.Stdout.Write(word)
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
	Count() int
	PermuteAll(letters []byte) [][]byte
}

type mtWrapper struct {
	trie  maptrie.Trie
	count int
}

func (m *mtWrapper) Load(word []byte) bool {
	m.trie = maptrie.Load(m.trie, word)
	m.count++
	return true
}

func (m *mtWrapper) PermuteAll(letters []byte) [][]byte {
	return maptrie.PermuteAll(m.trie, []byte(letters))
}

func (m *mtWrapper) Count() int {
	return m.count
}

type atWrapper struct {
	trie  *arraytrie.Trie
	count int
}

func (a *atWrapper) Load(word []byte) bool {
	a.trie = arraytrie.Load(a.trie, word)
	a.count++
	return true
}

func (a *atWrapper) PermuteAll(letters []byte) [][]byte {
	return arraytrie.PermuteAll(a.trie, []byte(letters))
}

func (a *atWrapper) Count() int {
	return a.count
}
