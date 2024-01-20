package main

import (
	"testing"
)

func TestNewRenumFolder(t *testing.T) {
	seasonNum := uint(2)
	epNum := uint(3)
	fileNames := []string{"S01E01.mkv", "S01E02.mkv"}

	processors := getProcessors()
	renumFolder := NewRenumFolder(seasonNum, epNum, fileNames)
	for _, renumFile := range renumFolder.RenumFiles {
		renumFile.Process(processors)
	}

	if renumFolder.SeasonNum != seasonNum {
		t.Errorf("Expected SeasonNum to be %d, but got %d", seasonNum, renumFolder.SeasonNum)
	}

	if renumFolder.EpNum != epNum {
		t.Errorf("Expected EpNum to be %d, but got %d", epNum, renumFolder.EpNum)
	}

	if len(renumFolder.RenumFiles) != len(fileNames) {
		t.Errorf("Expected RenumFiles length to be %d, but got %d", len(fileNames), len(renumFolder.RenumFiles))
	}

	if renumFolder.RenumFiles[0].NewName != "S02E03.mkv" {
		t.Errorf("Expected RenumFiles[0].NewName to be S02E03.mkv, but got %s", renumFolder.RenumFiles[0].NewName)
	}
	if renumFolder.RenumFiles[1].NewName != "S02E04.mkv" {
		t.Errorf("Expected RenumFiles[0].NewName to be S02E04.mkv, but got %s", renumFolder.RenumFiles[0].NewName)
	}
}
