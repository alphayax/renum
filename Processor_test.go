package main

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	p := NewProcessor("S[0-9]+E[0-9]+", "S%02dE%02d")
	if p.OutputPattern != "S%02dE%02d" {
		t.Errorf("Expected OutputPattern to be 'S%%02dE%%02d', but got %s", p.OutputPattern)
	}
	if !reflect.DeepEqual(p.SearchRegex, regexp.MustCompile("S[0-9]+E[0-9]+")) {
		t.Errorf("Expected SearchRegex to be 'S[0-9]+E[0-9]+', but got %v", p.SearchRegex)
	}
}

func TestMatch(t *testing.T) {
	p := NewProcessor("S[0-9]+E[0-9]+", "S%02dE%02d")
	if !p.Match("S01E01") {
		t.Errorf("Expected Match to return true, but got false")
	}
	if p.Match("S01F01") {
		t.Errorf("Expected Match to return false, but got true")
	}
}

func TestReplace(t *testing.T) {
	p := NewProcessor("S[0-9]+E[0-9]+", "S%02dE%02d")
	if p.Replace("S01E01", 2, 3) != "S02E03" {
		t.Errorf("Expected Replace to return 'S02E03', but got %s", p.Replace("S01E01", 2, 3))
	}
}
