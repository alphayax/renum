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
	if config.Json {
		log.SetFormatter(&log.JSONFormatter{})
	}
	if config.DryRun {
		log.Infoln("[DRY RUN] Dry run mode enabled, nothing will be changed")
	}

	// Print config
	log.WithFields(log.Fields{
		"SeasonNumber":       config.SeasonNum,
		"startEpisodeNumber": config.EpNum,
		"folder":             config.Folder,
		"dryRun":             config.DryRun,
		"force":              config.Force,
		"searchPattern":      config.SearchPattern,
	}).Debugln("[Config]")

	// Get the processors
	processors, err := getProcessors(config)
	if err != nil {
		log.Fatalln(err)
	}
	for _, processor := range processors {
		log.WithFields(log.Fields{
			"searchRegex":   processor.SearchRegex,
			"outputPattern": processor.OutputPattern,
		}).Debugln("[Processor]")
	}

	// Get the file names to process and compute the new names
	fileNames, err := getFolderFileNames(config.Folder)
	if err != nil {
		log.Fatalln(err)
	}
	renumFolder := NewRenumFolder(config.SeasonNum, config.EpNum, fileNames, processors)
	for _, file := range renumFolder.RenumFiles {
		log.WithFields(log.Fields{
			"oldName": file.OldName,
			"newName": file.NewName,
		}).Infoln("[Preview]")
	}

	// Refuse a plan that would overwrite files, before touching anything
	plan := buildRenamePlan(renumFolder.RenumFiles)
	if err := validateRenamePlan(config.Folder, plan); err != nil {
		log.Errorln("The following renames would destroy files:")
		for _, conflict := range strings.Split(err.Error(), "\n") {
			log.Errorln("  -", conflict)
		}
		log.Fatalln("Aborting, no file has been changed")
	}

	if config.DryRun {
		log.Infoln("[DRY RUN] Exiting...")
		os.Exit(0)
	}

	if len(plan) == 0 {
		log.Infoln("Nothing to rename")
		os.Exit(0)
	}

	// Ask for confirmation
	if !isOperationConfirmed(config.Force) {
		log.Warningln("Aborting the operation...")
		os.Exit(1)
	}

	// Rename files
	log.Infoln("Continuing the operation...")
	for _, op := range plan {
		log.WithFields(log.Fields{
			"oldName": op.From,
			"newName": op.To,
		}).Debugln("[Rename]")
	}
	if err := applyRenamePlan(config.Folder, plan); err != nil {
		for _, failure := range strings.Split(err.Error(), "\n") {
			log.Errorln(failure)
		}
		log.Fatalln("Some files could not be renamed")
	}

	log.Infof("Renamed %d file(s)", len(plan))
}

func getProcessors(config *Config) ([]*Processor, error) {
	if config.SearchPattern == "" {
		return getDefaultProcessors(), nil
	}

	processor, err := NewProcessor(config.SearchPattern, defaultOutputPattern)
	if err != nil {
		return nil, err
	}
	log.Infoln("Using custom search pattern:", config.SearchPattern)

	return []*Processor{processor}, nil
}

func isOperationConfirmed(force bool) bool {
	if force {
		log.Infoln("Force mode enabled, continuing the operation...")
		return true
	}

	fmt.Print("Do you want to continue the operation? (y/N): ")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		log.Debugln("Error while reading the response:", err)
		log.Debugln("Assuming the response is 'n'")
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y"
}

// getFolderFileNames lists the files of a folder, sub-folders excluded.
func getFolderFileNames(folderPath string) ([]string, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read the folder: %w", err)
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}
