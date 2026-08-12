package main

import (
	"testing"
)

// newNames maps each old name to the name it would be renamed to.
func newNames(folder *RenumFolder) map[string]string {
	names := make(map[string]string, len(folder.RenumFiles))
	for _, renumFile := range folder.RenumFiles {
		names[renumFile.OldName] = renumFile.NewName
	}

	return names
}

func TestNewRenumFolder(t *testing.T) {
	seasonNum := uint(2)
	epNum := uint(3)
	fileNames := []string{"S01E01.mkv", "S01E02.mkv"}

	renumFolder := NewRenumFolder(seasonNum, epNum, fileNames, getDefaultProcessors())

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
		t.Errorf("Expected RenumFiles[1].NewName to be S02E04.mkv, but got %s", renumFolder.RenumFiles[1].NewName)
	}
}

// A subtitle used to consume a number of its own, which desynchronised it from
// its video and shifted every episode that followed.
func TestNewRenumFolderKeepsCompanionFilesOnTheirEpisode(t *testing.T) {
	fileNames := []string{
		"Show S01E01.mkv", "Show S01E01.srt",
		"Show S01E02.mkv", "Show S01E02.srt",
	}

	names := newNames(NewRenumFolder(2, 1, fileNames, getDefaultProcessors()))
	expected := map[string]string{
		"Show S01E01.mkv": "Show S02E01.mkv",
		"Show S01E01.srt": "Show S02E01.srt",
		"Show S01E02.mkv": "Show S02E02.mkv",
		"Show S01E02.srt": "Show S02E02.srt",
	}
	for oldName, want := range expected {
		if names[oldName] != want {
			t.Errorf("Expected %q to become %q, but got %q", oldName, want, names[oldName])
		}
	}
}

// A file no processor matches is not an episode: numbering it would shift every
// episode coming after it.
func TestNewRenumFolderIgnoresFilesThatMatchNothing(t *testing.T) {
	fileNames := []string{"cover.jpg", "Show S01E01.mkv", "Show S01E02.mkv"}

	names := newNames(NewRenumFolder(1, 1, fileNames, getDefaultProcessors()))
	if names["cover.jpg"] != "cover.jpg" {
		t.Errorf("Expected cover.jpg to be left alone, but got %q", names["cover.jpg"])
	}
	if names["Show S01E01.mkv"] != "Show S01E01.mkv" {
		t.Errorf("Expected the first episode to keep number 1, but got %q", names["Show S01E01.mkv"])
	}
	if names["Show S01E02.mkv"] != "Show S01E02.mkv" {
		t.Errorf("Expected the second episode to keep number 2, but got %q", names["Show S01E02.mkv"])
	}
}

func TestGroupByEpisode(t *testing.T) {
	groups := groupByEpisode([]string{"a.mkv", "a.srt", "b.mkv"})
	if len(groups) != 2 {
		t.Fatalf("Expected 2 episodes, but got %d: %v", len(groups), groups)
	}
	if len(groups[0]) != 2 || groups[0][0] != "a.mkv" || groups[0][1] != "a.srt" {
		t.Errorf("Expected a.mkv and a.srt to be grouped, but got %v", groups[0])
	}
	if len(groups[1]) != 1 || groups[1][0] != "b.mkv" {
		t.Errorf("Expected b.mkv to be alone, but got %v", groups[1])
	}
}

func TestEpisodeKey(t *testing.T) {
	if key := episodeKey("Show S01E01.mkv"); key != "Show S01E01" {
		t.Errorf("Expected the extension to be dropped, but got %q", key)
	}
	if key := episodeKey("noextension"); key != "noextension" {
		t.Errorf("Expected the name to be kept as is, but got %q", key)
	}
}

func TestMatchesAny(t *testing.T) {
	processors := getDefaultProcessors()
	if !matchesAny([]string{"cover.jpg", "Show S01E01.mkv"}, processors) {
		t.Errorf("Expected the group to match through its episode file")
	}
	if matchesAny([]string{"cover.jpg"}, processors) {
		t.Errorf("Expected cover.jpg to match nothing")
	}
}
