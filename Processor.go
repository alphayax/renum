package main

import (
	"fmt"
	"regexp"
)

// defaultOutputPattern is the naming scheme applied when a file matches.
const defaultOutputPattern = "S%02dE%02d"

type Processor struct {
	OutputPattern string
	SearchRegex   *regexp.Regexp
}

// NewProcessor builds a Processor from a search pattern. It returns an error if
// the pattern is not a valid regular expression, which is why it must be used
// for any pattern coming from the user (see the --pattern flag).
func NewProcessor(searchPattern string, outputPattern string) (*Processor, error) {
	searchRegex, err := regexp.Compile(searchPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid search pattern %q: %w", searchPattern, err)
	}

	return &Processor{
		OutputPattern: outputPattern,
		SearchRegex:   searchRegex,
	}, nil
}

// MustNewProcessor is like NewProcessor but panics on an invalid pattern. It is
// reserved for the hard-coded default patterns, never for user input.
func MustNewProcessor(searchPattern string, outputPattern string) *Processor {
	processor, err := NewProcessor(searchPattern, outputPattern)
	if err != nil {
		panic(err)
	}

	return processor
}

func (p *Processor) Match(oldName string) bool {
	return p.SearchRegex.MatchString(oldName)
}

func (p *Processor) Replace(oldName string, seasonNum uint, epNum uint) string {
	replacement := fmt.Sprintf(p.OutputPattern, seasonNum, epNum)
	return p.SearchRegex.ReplaceAllString(oldName, replacement)
}

func getDefaultProcessors() []*Processor {
	return []*Processor{
		MustNewProcessor("S[0-9]+E[0-9]+", defaultOutputPattern),
		MustNewProcessor(" [0-9]{1,2}x[0-9]+ ", " "+defaultOutputPattern+" "),
		MustNewProcessor("^E[0-9]+", defaultOutputPattern),
		MustNewProcessor("[_ ][0-9]+[_ .]", "_"+defaultOutputPattern+"_"),
	}
}
