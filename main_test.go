package main

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// A sub-folder used to shift the whole listing: the slice was filled at the
// index of the directory entry, so an empty name appeared and the last file was
// dropped. The folder is named to sort before the files on purpose.
func TestGetFolderFileNamesSkipsSubFolders(t *testing.T) {
	folder := makeFolder(t, "Show S01E01.mkv", "Show S01E02.mkv", "Show S01E03.mkv")
	if err := os.Mkdir(filepath.Join(folder, "Extras"), 0o700); err != nil {
		t.Fatalf("unable to create the sub-folder: %v", err)
	}

	fileNames, err := getFolderFileNames(folder)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	expected := []string{"Show S01E01.mkv", "Show S01E02.mkv", "Show S01E03.mkv"}
	if len(fileNames) != len(expected) {
		t.Fatalf("Expected %d files, but got %d: %v", len(expected), len(fileNames), fileNames)
	}
	for i, want := range expected {
		if fileNames[i] != want {
			t.Errorf("Expected file %d to be %q, but got %q", i, want, fileNames[i])
		}
	}
}

// os.ReadDir hands back the lexicographic order, in which "ep 10.mkv" precedes
// "ep 2.mkv" and every episode gets the number of another.
func TestGetFolderFileNamesUsesTheNaturalOrder(t *testing.T) {
	folder := makeFolder(t, "ep 10.mkv", "ep 2.mkv", "ep 1.mkv")

	fileNames, err := getFolderFileNames(folder)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	expected := []string{"ep 1.mkv", "ep 2.mkv", "ep 10.mkv"}
	if !slices.Equal(fileNames, expected) {
		t.Errorf("Expected %v, but got %v", expected, fileNames)
	}
}

func TestGetFolderFileNamesOnMissingFolder(t *testing.T) {
	if _, err := getFolderFileNames(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Errorf("Expected an error on a missing folder, but got nil")
	}
}

func TestGetProcessorsUsesTheDefaults(t *testing.T) {
	processors, err := getProcessors(&Config{})
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if len(processors) != len(getDefaultProcessors()) {
		t.Errorf("Expected the default processors, but got %d of them", len(processors))
	}
}

func TestGetProcessorsUsesTheCustomPattern(t *testing.T) {
	processors, err := getProcessors(&Config{SearchPattern: "Ep[0-9]+"})
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if len(processors) != 1 {
		t.Fatalf("Expected 1 processor, but got %d", len(processors))
	}
	if !processors[0].Match("Ep12") {
		t.Errorf("Expected the custom pattern to match 'Ep12'")
	}
}

func TestGetProcessorsRejectsAnInvalidPattern(t *testing.T) {
	if _, err := getProcessors(&Config{SearchPattern: "S[0-9"}); err == nil {
		t.Errorf("Expected an invalid pattern to return an error")
	}
}

// The last gate before anything is renamed: anything but a plain "y" — an
// answer that cannot even be read included — must decline.
func TestIsOperationConfirmed(t *testing.T) {
	cases := map[string]struct {
		force bool
		input string
		want  bool
	}{
		"force":         {force: true, input: "", want: true},
		"y":             {input: "y\n", want: true},
		"uppercase Y":   {input: "Y\n", want: true},
		"padded y":      {input: "  y  \n", want: true},
		"n":             {input: "n\n", want: false},
		"empty line":    {input: "\n", want: false},
		"closed input":  {input: "", want: false},
		"anything else": {input: "yes\n", want: false},
	}

	for name, testCase := range cases {
		got := isOperationConfirmed(testCase.force, strings.NewReader(testCase.input))
		if got != testCase.want {
			t.Errorf("%s: expected %t, but got %t", name, testCase.want, got)
		}
	}
}

// errors.Join glues its errors with a newline, and a file name may hold one:
// the messages are taken apart, never split back on the newline.
func TestJoinedErrors(t *testing.T) {
	first := errors.New("first")
	second := errors.New("a name\nwith a newline")

	joined := joinedErrors(errors.Join(first, second))
	if len(joined) != 2 {
		t.Fatalf("Expected 2 errors, but got %d: %v", len(joined), joined)
	}
	if joined[0] != first || joined[1] != second {
		t.Errorf("Expected the errors as they were joined, but got %v", joined)
	}

	alone := joinedErrors(first)
	if len(alone) != 1 || alone[0] != first {
		t.Errorf("Expected a lone error to come back as it is, but got %v", alone)
	}
}
