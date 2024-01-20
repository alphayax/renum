package main

import (
	"fmt"
	"regexp"
)

type Processor struct {
	OutputPattern string
	SearchRegex   *regexp.Regexp
}

func NewProcessor(searchPattern string, outputPattern string) *Processor {
	return &Processor{
		OutputPattern: outputPattern,
		SearchRegex:   regexp.MustCompile(searchPattern),
	}
}

func (p *Processor) Match(oldName string) bool {
	return p.SearchRegex.MatchString(oldName)
}

func (p *Processor) Replace(oldName string, seasonNum uint, epNum uint) string {
	replacement := fmt.Sprintf(p.OutputPattern, seasonNum, epNum)
	return p.SearchRegex.ReplaceAllString(oldName, replacement)
}

func getProcessors() []*Processor {
	return []*Processor{
		NewProcessor("S[0-9]+E[0-9]+", "S%02dE%02d"),
		NewProcessor(" [0-9]{1,2}x[0-9]+ ", " S%02dE%02d "),
		NewProcessor("^E[0-9]+", "S%02dE%02d"),
		NewProcessor("[_ ][0-9]+[_ .]", "_S%02dE%02d_"),
	}
}
