package main

import (
	"fmt"
	"regexp"
)

func NewRenumFile(oldName string, seasonNum uint, epNum uint) *RenumFile {
	renumFile := &RenumFile{
		OldName:   oldName,
		SeasonNum: seasonNum,
		EpNum:     epNum,
	}
	renumFile.NewName = renumFile.getNewName()
	return renumFile
}

type RenumFile struct {
	OldName   string
	NewName   string
	SeasonNum uint
	EpNum     uint
}

func (r *RenumFile) getNewName() string {

	// Search for something like S01E01
	if re := regexp.MustCompile("S[0-9]+E[0-9]+"); re.MatchString(r.OldName) {
		//log.Println("Match 'S[0-9]+E[0-9]+'")
		return re.ReplaceAllString(r.OldName, fmt.Sprintf("S%02dE%02d", r.SeasonNum, r.EpNum))
	}

	// Search for something like 01x01 or 1x01
	if re := regexp.MustCompile(" [0-9]{1,2}x[0-9]+ "); re.MatchString(r.OldName) {
		//log.Println("Match ' [0-9]{1,2}x[0-9]+ '")
		return re.ReplaceAllString(r.OldName, fmt.Sprintf(" S%02dE%02d ", r.SeasonNum, r.EpNum))
	}

	// Search for something who start by E01 or E1
	if re := regexp.MustCompile("^E[0-9]+"); re.MatchString(r.OldName) {
		//log.Println("Match '^E[0-9]+'")
		return re.ReplaceAllString(r.OldName, fmt.Sprintf("S%02dE%02d", r.SeasonNum, r.EpNum))
	}

	// Search for something like _001_ or _01_ or 001.
	if re := regexp.MustCompile("[_ ][0-9]+[_ .]"); re.MatchString(r.OldName) {
		//log.Println("Match '[_ ][0-9]+[_ .]'")
		return re.ReplaceAllString(r.OldName, fmt.Sprintf("_S%02dE%02d_", r.SeasonNum, r.EpNum))
	}

	//log.Println("No match")
	return r.OldName
}

func (r *RenumFile) String() string {
	return fmt.Sprintf("%s -> %s", r.OldName, r.NewName)
}
