package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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
	// Parsing
	flags  *flag.FlagSet
	output io.Writer
}

func NewConfig() *Config {
	return newConfig(os.Stderr)
}

// newConfig builds a Config on a flag set of its own, and makes it write what it
// has to say to output. The global flag set the options used to be declared on
// panics the second time a Config is built in a process, which is exactly what
// a test needs to do.
func newConfig(output io.Writer) *Config {
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
		flags:         flag.NewFlagSet("renum", flag.ContinueOnError),
		output:        output,
	}

	// renum reports the parsing errors and prints the usage itself, once, from
	// main: the flag set is told to keep quiet.
	c.flags.SetOutput(io.Discard)

	c.flags.UintVar(&c.SeasonNum, "season", 1, "New season number")
	c.flags.UintVar(&c.SeasonNum, "s", 1, "New season number")
	c.flags.UintVar(&c.EpNum, "episode", 1, "Starting episode number")
	c.flags.UintVar(&c.EpNum, "e", 1, "Starting episode number")
	c.flags.BoolVar(&c.Verbose, "verbose", false, "Increase verbosity")
	c.flags.BoolVar(&c.Json, "json", false, "Set logs into json")
	c.flags.BoolVar(&c.DryRun, "dry-run", false, "Does nothing, just print the new names")
	c.flags.BoolVar(&c.Help, "h", false, "Print this help message")
	c.flags.BoolVar(&c.Help, "help", false, "Print this help message")
	c.flags.BoolVar(&c.Force, "force", false, "Don't ask for confirmation")
	c.flags.StringVar(&c.SearchPattern, "pattern", "", "Custom regex search pattern (eg: S[0-9]+E[0-9]+)")

	return c
}

func (c *Config) Parse(args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return err
	}

	if c.flags.NArg() > 0 {
		c.Folder = c.flags.Arg(0)
	}

	// Asking for the help is not a failure: the folder nobody passes with -h
	// must not turn it into one.
	if c.Help {
		return nil
	}

	if c.flags.NArg() != 1 {
		return fmt.Errorf("invalid number of arguments (got %d, expected 1)", c.flags.NArg())
	}

	return nil
}

func (c *Config) Usage() {
	// Printing the defaults is the one moment the flag set is given the output
	// back.
	c.flags.SetOutput(c.output)
	defer c.flags.SetOutput(io.Discard)

	_, _ = fmt.Fprintln(c.output, "Usage: renum [options] <folderPath>")
	_, _ = fmt.Fprintln(c.output, "Options:")
	c.flags.PrintDefaults()
}
