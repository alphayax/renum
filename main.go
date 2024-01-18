package main

import (
	"fmt"
	"log"
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
	if config.DryRun {
		log.Println("[DRY RUN] Dry run mode enabled, nothing will be changed")
	}

	fileNames := getFolderFileNames(config.Folder)
	renumFolder := NewRenumFolder(config.SeasonNum, config.EpNum, fileNames)
	for _, file := range renumFolder.RenumFiles {
		log.Printf("[Preview] %s\n", file.String())
	}

	if config.DryRun {
		log.Println("[DRY RUN] Exiting...")
		os.Exit(0)
	}

	// Ask for confirmation
	if !isOperationConfirmed(config.Force) {
		log.Println("Aborting the operation...")
		os.Exit(-1)
	}

	log.Println("Continuing the operation...")
	for _, file := range renumFolder.RenumFiles {
		log.Printf("[Rename] %s\n", file.String())
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
		log.Fatal(err)
	}

	// TODO: Exclude directories
	fileName := make([]string, len(files))
	for i, file := range files {
		fileName[i] = file.Name()
	}
	return fileName
}
