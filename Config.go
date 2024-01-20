package main

import (
	"flag"
	"fmt"
)

type Config struct {
	SeasonNum uint
	EpNum     uint
	Verbose   bool
	DryRun    bool
	Help      bool
	Force     bool
	Folder    string
}

func NewConfig() *Config {
	c := &Config{
		SeasonNum: 1,
		EpNum:     1,
		Verbose:   false,
		DryRun:    false,
		Help:      false,
		Force:     false,
		Folder:    "",
	}

	flag.UintVar(&c.SeasonNum, "season", 1, "New season number")
	flag.UintVar(&c.SeasonNum, "s", 1, "New season number")
	flag.UintVar(&c.EpNum, "episode", 1, "Starting episode number")
	flag.UintVar(&c.EpNum, "e", 1, "Starting episode number")
	flag.BoolVar(&c.Verbose, "verbose", false, "Increase verbosity")
	flag.BoolVar(&c.DryRun, "dry-run", false, "Does nothing, just print the new names")
	flag.BoolVar(&c.Help, "h", false, "Print this help message")
	flag.BoolVar(&c.Help, "help", false, "Print this help message")
	flag.BoolVar(&c.Force, "force", false, "Don't ask for confirmation")
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
