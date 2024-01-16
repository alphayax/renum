package main

import (
	"testing"
)

func TestNewRenumFile(t *testing.T) {
	oldName := "XXXX [ABCD] - S01E01 - Lorem Ipsum [1280x1024].mkv"
	seasonNum := uint(2)
	epNum := uint(3)
	renumFile := NewRenumFile(oldName, seasonNum, epNum)

	if renumFile.OldName != oldName {
		t.Errorf("Expected OldName to be %s, but got %s", oldName, renumFile.OldName)
	}

	if renumFile.SeasonNum != seasonNum {
		t.Errorf("Expected SeasonNum to be %d, but got %d", seasonNum, renumFile.SeasonNum)
	}

	if renumFile.EpNum != epNum {
		t.Errorf("Expected EpNum to be %d, but got %d", epNum, renumFile.EpNum)
	}
}

func TestGetNewName(t *testing.T) {
	oldName := "XXXX [ABCD] - S01E01 - Lorem Ipsum [1280x1024].mkv"
	seasonNum := uint(2)
	epNum := uint(3)
	renumFile := NewRenumFile(oldName, seasonNum, epNum)

	expectedNewName := "XXXX [ABCD] - S02E03 - Lorem Ipsum [1280x1024].mkv"
	newName := renumFile.getNewName()

	if newName != expectedNewName {
		t.Errorf("Expected NewName to be %s, but got %s", expectedNewName, newName)
	}
}
