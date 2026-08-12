package main

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	p, err := NewProcessor("S[0-9]+E[0-9]+", "S%02dE%02d")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if p.OutputPattern != "S%02dE%02d" {
		t.Errorf("Expected OutputPattern to be 'S%%02dE%%02d', but got %s", p.OutputPattern)
	}
	if !reflect.DeepEqual(p.SearchRegex, regexp.MustCompile("S[0-9]+E[0-9]+")) {
		t.Errorf("Expected SearchRegex to be 'S[0-9]+E[0-9]+', but got %v", p.SearchRegex)
	}
}

// An invalid --pattern must be reported as an error, never as a panic.
func TestNewProcessorInvalidPattern(t *testing.T) {
	p, err := NewProcessor("S[0-9", "S%02dE%02d")
	if err == nil {
		t.Fatalf("Expected an error for an invalid pattern, but got processor %v", p)
	}
	if p != nil {
		t.Errorf("Expected a nil processor on error, but got %v", p)
	}
}

func TestMustNewProcessor(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Errorf("Expected MustNewProcessor to panic on an invalid pattern")
		}
	}()
	MustNewProcessor("S[0-9", "S%02dE%02d")
}

func TestMatch(t *testing.T) {
	p := MustNewProcessor("S[0-9]+E[0-9]+", "S%02dE%02d")
	if !p.Match("S01E01") {
		t.Errorf("Expected Match to return true, but got false")
	}
	if p.Match("S01F01") {
		t.Errorf("Expected Match to return false, but got true")
	}
}

func TestReplace(t *testing.T) {
	p := MustNewProcessor("S[0-9]+E[0-9]+", "S%02dE%02d")
	if p.Replace("S01E01", 2, 3) != "S02E03" {
		t.Errorf("Expected Replace to return 'S02E03', but got %s", p.Replace("S01E01", 2, 3))
	}
}

func TestGetDefaultProcessors(t *testing.T) {
	if len(getDefaultProcessors()) == 0 {
		t.Errorf("Expected the default processors not to be empty")
	}
}

// The separator pattern used to swallow the dot of the extension, turning
// "serie 1.mkv" into "serie_S01E01_mkv".
func TestDefaultProcessorsKeepTheExtension(t *testing.T) {
	cases := map[string]string{
		"serie 1.mkv":                     "serie S01E01.mkv",
		"serie_1.mkv":                     "serie_S01E01.mkv",
		"serie 1 x.mkv":                   "serie S01E01 x.mkv",
		"[Fansub]_Show_1086_[VOSTFR].mkv": "[Fansub]_Show_S01E01_[VOSTFR].mkv",
	}

	for oldName, want := range cases {
		names := newNames(NewRenumFolder(1, 1, []string{oldName}, getDefaultProcessors()))
		if names[oldName] != want {
			t.Errorf("Expected %q to become %q, but got %q", oldName, want, names[oldName])
		}
	}
}
