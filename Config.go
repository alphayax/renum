package main

import (
	"flag"
	"fmt"
)

type Config struct {
	// Options
	SeasonNum     uint
	EpNum         uint
	Verbose       bool
	Json          bool
	DryRun        bool
	Help          bool
	Force         bool
	SearchPattern string
	// Positional arguments
	Folder string
}

func NewConfig() *Config {
	c := &Config{
		SeasonNum:     1,
		EpNum:         1,
		Verbose:       false,
		DryRun:        false,
		Help:          false,
		Force:         false,
		Json:          false,
		SearchPattern: "",
		Folder:        "",
	}

	flag.UintVar(&c.SeasonNum, "season", 1, "New season number")
	flag.UintVar(&c.SeasonNum, "s", 1, "New season number")
	flag.UintVar(&c.EpNum, "episode", 1, "Starting episode number")
	flag.UintVar(&c.EpNum, "e", 1, "Starting episode number")
	flag.BoolVar(&c.Verbose, "verbose", false, "Increase verbosity")
	flag.BoolVar(&c.Json, "json", false, "Set logs into json")
	flag.BoolVar(&c.DryRun, "dry-run", false, "Does nothing, just print the new names")
	flag.BoolVar(&c.Help, "h", false, "Print this help message")
	flag.BoolVar(&c.Help, "help", false, "Print this help message")
	flag.BoolVar(&c.Force, "force", false, "Don't ask for confirmation")
	flag.StringVar(&c.SearchPattern, "pattern", "", "Custom regex search pattern (eg: S[0-9]+E[0-9]+)")
	flag.Usage = func() {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Usage: renum [options] <folderPath>")
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}

	return c
}

func (c *Config) Parse() error {
	flag.Parse()
	if flag.NArg() != 1 {
		return fmt.Errorf("invalid number of arguments (got %d, expected 1)", flag.NArg())
	}
	c.Folder = flag.Arg(0)
	return nil
}

func (c *Config) Usage() {
	flag.Usage()
}
