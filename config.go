package main

import (
	"errors"

	"github.com/spf13/pflag"
)

type Config struct {
	UseArray bool
	Dict     string
	Verbose  bool
	Consume  bool
	All      bool
}

func AddFlags(flags *pflag.FlagSet) {
	flags.Bool("use-array", false, "Use the arrayTrie implementation instead of the mapTrie")
	flags.StringP("dict", "d", "2of12inf", "Name of the dictionary to use")
	flags.BoolP("verbose", "v", false, "Get wordy with those words")
	flags.BoolP("consume", "c", false, "Consume letters (only use each letter once)")
	flags.BoolP("all", "a", false, "Use all letters")
}

func GetConfig(flags *pflag.FlagSet) (Config, error) {
	var (
		// getint    = (*pflag.FlagSet).GetInt
		getstring = (*pflag.FlagSet).GetString
		getbool   = (*pflag.FlagSet).GetBool
	)

	var cfg Config
	return cfg, errors.Join(
		GetFlagT(&cfg.UseArray, flags, "use-array", getbool),
		GetFlagT(&cfg.Dict, flags, "dict", getstring),
		GetFlagT(&cfg.Verbose, flags, "verbose", getbool),
		GetFlagT(&cfg.Consume, flags, "consume", getbool),
		GetFlagT(&cfg.All, flags, "all", getbool),
	)
}

func GetFlagT[T any](target *T, flags *pflag.FlagSet, field string, extractor func(*pflag.FlagSet, string) (T, error)) error {
	var err error
	*target, err = extractor(flags, field)
	return err
}
