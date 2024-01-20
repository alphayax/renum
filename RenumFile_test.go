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

var renumFileTestData = []struct {
	in  string
	out string
}{
	{"XXXX [ABCD] - S01E01 - Lorem Ipsum [1280x1024].mkv", "XXXX [ABCD] - S02E03 - Lorem Ipsum [1280x1024].mkv"},
	{"XXXX [ABCD] - 01x01 - Lorem Ipsum [1280x1024].mkv", "XXXX [ABCD] - S02E03 - Lorem Ipsum [1280x1024].mkv"},
	{"XXXX [ABCD] - 1x01 - Lorem Ipsum [1280x1024].mkv", "XXXX [ABCD] - S02E03 - Lorem Ipsum [1280x1024].mkv"},
	{"XXXX_[ABCD]_001_Lorem_Ipsum_[1280x1024].mkv", "XXXX_[ABCD]_S02E03_Lorem_Ipsum_[1280x1024].mkv"},
	{"E01 - Lorem Ipsum [1280x1024].mkv", "S02E03 - Lorem Ipsum [1280x1024].mkv"},
	{"ABCDEF.mkv", "ABCDEF.mkv"},
}

func TestGetNewName(t *testing.T) {
	for _, tt := range renumFileTestData {
		processors := getProcessors()
		t.Run(tt.in, func(t *testing.T) {
			renumFile := NewRenumFile(tt.in, 2, 3)
			if renumFile.getNewName(processors) != tt.out {
				t.Errorf("Expected NewName to be %s, but got %s", tt.out, renumFile.getNewName(processors))
			}
		})
	}
}

func TestString(t *testing.T) {
	for _, tt := range renumFileTestData {
		processors := getProcessors()
		t.Run(tt.in, func(t *testing.T) {
			renumFile := NewRenumFile(tt.in, 2, 3)
			renumFile.Process(processors)
			if renumFile.String() != tt.in+" -> "+tt.out {
				t.Errorf("Expected String to be %s, but got %s", tt.in+" -> "+tt.out, renumFile.String())
			}
		})
	}
}
