package main

import "log"

type RenumFolder struct {
	SeasonNum  uint
	EpNum      uint
	FolderPath string
	RenumFiles []*RenumFile
}

func NewRenumFolder(seasonNum uint, epNum uint, folder string, fileNames []string) *RenumFolder {
	renumFolder := &RenumFolder{
		SeasonNum:  seasonNum,
		EpNum:      epNum,
		FolderPath: folder,
		RenumFiles: make([]*RenumFile, len(fileNames)),
	}

	var i uint = 0
	for _, fileName := range fileNames {
		renumFolder.RenumFiles[i] = NewRenumFile(fileName, seasonNum, epNum+i)
		i++
	}

	return renumFolder
}

func (r *RenumFolder) Preview() {
	for _, file := range r.RenumFiles {
		file.Preview()
	}
}

func (r *RenumFolder) Rename() {
	for _, file := range r.RenumFiles {
		if err := file.Rename(r.FolderPath); err != nil {
			log.Fatal(err)
		}
	}
}
