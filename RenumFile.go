package main

func NewRenumFile(oldName string, seasonNum uint, epNum uint) *RenumFile {
	return &RenumFile{
		OldName:   oldName,
		SeasonNum: seasonNum,
		EpNum:     epNum,
	}
}

type RenumFile struct {
	OldName   string
	NewName   string
	SeasonNum uint
	EpNum     uint
}

func (r *RenumFile) Process(processors []*Processor) {
	r.NewName = r.getNewName(processors)
}

func (r *RenumFile) getNewName(processors []*Processor) string {
	for _, processor := range processors {
		if processor.Match(r.OldName) {
			return processor.Replace(r.OldName, r.SeasonNum, r.EpNum)
		}
	}

	return r.OldName
}
