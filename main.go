package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	seasonNum := flag.Int("season", 1, "New season number")
	epNum := flag.Int("episode", 1, "Starting episode number")
	dryRun := flag.Bool("dry-run", false, "Does nothing, just print the new names")
	flag.Parse()
	NArg := flag.NArg()
	folder := flag.Arg(0)
	flag.Usage = func() {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Usage: renum [options] <folder>")
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}

	if NArg != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if dryRun != nil && *dryRun {
		log.Println("[DRY RUN] Dry run mode enabled, nothing will be changed")
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	renumFiles := make([]*RenumFile, len(files))
	i := 0
	for _, file := range files {
		renumFiles[i] = NewRenumFile(file.Name(), uint(*seasonNum), uint(*epNum+i))
		renumFiles[i].Preview()
		i++
	}

	if dryRun != nil && *dryRun {
		log.Println("[DRY RUN] Exiting...")
		os.Exit(0)
	}

	// Ask for confirmation
	if !isOperationConfirmed() {
		log.Println("Aborting the operation...")
		os.Exit(-1)
	}

	log.Println("Continuing the operation...")
	for _, renumFile := range renumFiles {
		if err := renumFile.Rename(folder); err != nil {
			log.Fatal(err)
		}
	}
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
