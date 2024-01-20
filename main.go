package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	config := NewConfig()
	if err := config.Parse(); err != nil {
		fmt.Println(err)
		config.Usage()
		os.Exit(1)
	}
	if config.Help {
		config.Usage()
		os.Exit(0)
	}
	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	if config.DryRun {
		log.Infoln("[DRY RUN] Dry run mode enabled, nothing will be changed")
	}

	log.WithFields(log.Fields{
		"SeasonNumber":       config.SeasonNum,
		"startEpisodeNumber": config.EpNum,
		"folder":             config.Folder,
		"dryRun":             config.DryRun,
	}).Debugln("[Config]")

	fileNames := getFolderFileNames(config.Folder)
	renumFolder := NewRenumFolder(config.SeasonNum, config.EpNum, fileNames)

	processors := getProcessors()
	for _, processor := range processors {
		log.WithFields(log.Fields{
			"searchRegex":   processor.SearchRegex,
			"outputPattern": processor.OutputPattern,
		}).Debugln("[Processor]")
	}
	for _, file := range renumFolder.RenumFiles {
		file.Process(processors)
		log.WithFields(log.Fields{
			"oldName": file.OldName,
			"newName": file.NewName,
		}).Infoln("[Preview]")
	}

	if config.DryRun {
		log.Infoln("[DRY RUN] Exiting...")
		os.Exit(0)
	}

	// Ask for confirmation
	if !isOperationConfirmed(config.Force) {
		log.Warningln("Aborting the operation...")
		os.Exit(-1)
	}

	log.Infoln("Continuing the operation...")
	for _, file := range renumFolder.RenumFiles {
		log.WithFields(log.Fields{
			"oldName": file.OldName,
			"newName": file.NewName,
		}).Debugln("[Rename]")
		if err := os.Rename(
			fmt.Sprintf("%s/%s", config.Folder, file.OldName),
			fmt.Sprintf("%s/%s", config.Folder, file.NewName),
		); err != nil {
			log.Fatal(err)
		}
	}
}

func isOperationConfirmed(force bool) bool {
	if force {
		log.Infoln("Force mode enabled, continuing the operation...")
		return true
	}

	fmt.Print("Do you want to continue the operation? (y/N): ")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y"
}

func getFolderFileNames(folderPath string) []string {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalln(err)
	}

	fileNames := make([]string, len(files))
	folderCount := 0
	for i, file := range files {
		if file.IsDir() {
			folderCount++
			continue
		}
		fileNames[i] = file.Name()
	}

	// Resize the slice to remove the skipped folders
	fileNames = fileNames[:len(files)-folderCount]

	return fileNames
}
