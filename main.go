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
	var dryRun, help, force bool
	flag.UintVar(&seasonNum, "season", 1, "New season number")
	flag.UintVar(&seasonNum, "s", 1, "New season number")
	flag.UintVar(&epNum, "episode", 1, "Starting episode number")
	flag.UintVar(&epNum, "e", 1, "Starting episode number")
	flag.BoolVar(&dryRun, "dry-run", false, "Does nothing, just print the new names")
	flag.BoolVar(&help, "h", false, "Print this help message")
	flag.BoolVar(&help, "help", false, "Print this help message")
	flag.BoolVar(&force, "force", false, "Don't ask for confirmation")
	flag.Parse()
	flag.Usage = func() {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Usage: renum [options] <folderPath>")
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}

	if flag.NArg() != 1 || help {
		flag.Usage()
		os.Exit(1)
	}

	if dryRun {
		log.Println("[DRY RUN] Dry run mode enabled, nothing will be changed")
	}

	folderPath := flag.Arg(0)
	fileNames := getFolderFileNames(folderPath)
	renumFolder := NewRenumFolder(seasonNum, epNum, folderPath, fileNames)
	renumFolder.Preview()

	if dryRun {
		log.Println("[DRY RUN] Exiting...")
		os.Exit(0)
	}

	// Ask for confirmation
	if !isOperationConfirmed(force) {
		log.Println("Aborting the operation...")
		os.Exit(-1)
	}

	log.Println("Continuing the operation...")
	renumFolder.Rename()
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
