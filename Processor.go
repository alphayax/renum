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

// Replace renumbers the first occurrence of the pattern and leaves the rest of
// the name untouched. A name carries at most one episode number; any further
// occurrence belongs to the title, as in "S01E01 - rerun of S01E01.mkv", and
// renumbering it too would corrupt that title.
func (p *Processor) Replace(oldName string, seasonNum uint, epNum uint) string {
	match := p.SearchRegex.FindStringSubmatchIndex(oldName)
	if match == nil {
		return oldName
	}

	replacement := fmt.Sprintf(p.OutputPattern, seasonNum, epNum)
	// ExpandString gives the replacement the same "$1" semantics it had with
	// ReplaceAllString, which the separator pattern below relies on.
	expanded := p.SearchRegex.ExpandString(nil, replacement, oldName, match)

	return oldName[:match[0]] + string(expanded) + oldName[match[1]:]
}

func getDefaultProcessors() []*Processor {
	return []*Processor{
		MustNewProcessor("S[0-9]+E[0-9]+", defaultOutputPattern),
		MustNewProcessor(" [0-9]{1,2}x[0-9]+ ", " "+defaultOutputPattern+" "),
		MustNewProcessor("^E[0-9]+", defaultOutputPattern),
		// The surrounding separators are captured and put back as they were:
		// replacing them by underscores turned "serie 1.mkv" into
		// "serie_S01E01_mkv", stripping the extension off the file.
		MustNewProcessor("([_ ])[0-9]+([_ .])", "${1}"+defaultOutputPattern+"${2}"),
	}
}
