package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var seasonNum, epNum uint
	var dryRun, help bool
	flag.UintVar(&seasonNum, "season", 1, "New season number")
	flag.UintVar(&seasonNum, "s", 1, "New season number")
	flag.UintVar(&epNum, "episode", 1, "Starting episode number")
	flag.UintVar(&epNum, "e", 1, "Starting episode number")
	flag.BoolVar(&dryRun, "dry-run", false, "Does nothing, just print the new names")
	flag.BoolVar(&help, "h", false, "Print this help message")
	flag.BoolVar(&help, "help", false, "Print this help message")
	flag.Parse()
	NArg := flag.NArg()
	folder := flag.Arg(0)
	flag.Usage = func() {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Usage: renum [options] <folder>")
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}

	if NArg != 1 || help {
		flag.Usage()
		os.Exit(1)
	}

	if dryRun {
		log.Println("[DRY RUN] Dry run mode enabled, nothing will be changed")
	}

	fileNames := getFolderFileNames(folder)

	renumFolder := NewRenumFolder(seasonNum, epNum, folder, fileNames)
	renumFolder.Preview()

	if dryRun {
		log.Println("[DRY RUN] Exiting...")
		os.Exit(0)
	}

	// Ask for confirmation
	if !isOperationConfirmed() {
		log.Println("Aborting the operation...")
		os.Exit(-1)
	}

	log.Println("Continuing the operation...")
	renumFolder.Rename()
}

func isOperationConfirmed() bool {
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
