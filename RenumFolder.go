package main

type RenumFolder struct {
	SeasonNum  uint
	EpNum      uint
	FolderPath string
	RenumFiles []*RenumFile
}

func NewRenumFolder(seasonNum uint, epNum uint, fileNames []string) *RenumFolder {
	renumFolder := &RenumFolder{
		SeasonNum:  seasonNum,
		EpNum:      epNum,
		RenumFiles: make([]*RenumFile, len(fileNames)),
	}

	var i uint = 0
	for _, fileName := range fileNames {
		renumFolder.RenumFiles[i] = NewRenumFile(fileName, seasonNum, epNum+i)
		i++
	}

	return renumFolder
}
