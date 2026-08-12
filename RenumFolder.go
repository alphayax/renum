package main

import (
	"path/filepath"
	"strings"
)

type RenumFolder struct {
	SeasonNum  uint
	EpNum      uint
	RenumFiles []*RenumFile
}

// NewRenumFolder splits the files into episodes and gives each episode its
// number. All the files of one episode — a video and its subtitle, for
// instance — share that number, and a file that no processor matches is left
// out of the numbering entirely.
func NewRenumFolder(seasonNum uint, epNum uint, fileNames []string, processors []*Processor) *RenumFolder {
	renumFolder := &RenumFolder{
		SeasonNum:  seasonNum,
		EpNum:      epNum,
		RenumFiles: make([]*RenumFile, 0, len(fileNames)),
	}

	currentEpNum := epNum
	for _, group := range groupByEpisode(fileNames) {
		isEpisode := matchesAny(group, processors)
		for _, fileName := range group {
			renumFile := NewRenumFile(fileName, seasonNum, currentEpNum)
			renumFile.Process(processors)
			renumFolder.RenumFiles = append(renumFolder.RenumFiles, renumFile)
		}
		// Anything no processor matches is not an episode. Numbering it would
		// consume a number and shift every episode that follows.
		if isEpisode {
			currentEpNum++
		}
	}

	return renumFolder
}

// episodeKey is what files are grouped by: the name without its extension, so
// that a video and its subtitle end up in the same episode.
func episodeKey(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// groupByEpisode gathers the files of a same episode, keeping the order in
// which each episode first shows up.
func groupByEpisode(fileNames []string) [][]string {
	var groups [][]string

	positions := make(map[string]int, len(fileNames))
	for _, fileName := range fileNames {
		key := episodeKey(fileName)
		if position, seen := positions[key]; seen {
			groups[position] = append(groups[position], fileName)
			continue
		}
		positions[key] = len(groups)
		groups = append(groups, []string{fileName})
	}

	return groups
}

// matchesAny reports whether at least one file of the group is something Renum
// knows how to renumber.
func matchesAny(group []string, processors []*Processor) bool {
	for _, fileName := range group {
		for _, processor := range processors {
			if processor.Match(fileName) {
				return true
			}
		}
	}

	return false
}
