package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	seasonNum := flag.Int("season", 1, "New season number")
	dryRun := flag.Bool("dry-run", true, "Does nothing, just print the new names")
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
		fmt.Println("[DRY RUN] Dry run mode enabled, nothing will be changed")
	}

	//fmt.Println("seasonNum:", *seasonNum)

	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	epNum := 0
	for _, file := range files {
		epNum = epNum + 1

		newName := getNewName(file.Name(), *seasonNum, epNum)
		fmt.Println(file.Name() + "\t-> " + newName)

		if dryRun != nil && *dryRun {
			continue
		}

		err := os.Rename(folder+"/"+file.Name(), folder+"/"+newName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// getNewName returns the new name of the file with the season and the new episode number
func getNewName(oldName string, seasonNum int, epNum int) string {

	// Search for something like S01E01
	re := regexp.MustCompile("S[0-9]+E[0-9]+")
	if re.MatchString(oldName) {
		return re.ReplaceAllString(oldName, fmt.Sprintf("S%02dE%02d", seasonNum, epNum))
	}

	// Search for something like 01x01 or 1x01
	re = regexp.MustCompile("[0-9]{1-2}x[0-9]+")
	if re.MatchString(oldName) {
		return re.ReplaceAllString(oldName, fmt.Sprintf("S%02dE%02d", seasonNum, epNum))
	}

	// Search for something like _001_ or _01_ or 001.
	re = regexp.MustCompile("[_ ][0-9]+[_ .]")
	if re.MatchString(oldName) {
		return re.ReplaceAllString(oldName, fmt.Sprintf("_S%02dE%02d_", seasonNum, epNum))
	}

	return oldName
}
