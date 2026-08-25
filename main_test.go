package main

import (
	"os"
	"path/filepath"
	"slices"
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
