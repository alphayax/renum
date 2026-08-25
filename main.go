package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func main() {
	config := NewConfig()
	if err := config.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		config.Usage()
		os.Exit(1)
	}
	if config.Help {
		// Asking for the help is not a failure: it goes to the standard output,
		// so that "renum --help | less" and "renum --help > help.txt" see it.
		config.SetOutput(os.Stdout)
		config.Usage()
		os.Exit(0)
	}
	setupLogger(config)
	if config.DryRun {
		slog.Info("[DRY RUN] Dry run mode enabled, nothing will be changed")
	}

	// Print config
	slog.Debug("[Config]",
		"SeasonNumber", config.SeasonNum,
		"startEpisodeNumber", config.EpNum,
		"folder", config.Folder,
		"dryRun", config.DryRun,
		"force", config.Force,
		"searchPattern", config.SearchPattern,
	)

	// Get the processors
	processors, err := getProcessors(config)
	if err != nil {
		fatal(err.Error())
	}
	for _, processor := range processors {
		slog.Debug("[Processor]",
			// The regexp itself would be marshalled as an empty object in the
			// JSON output: what is wanted here is the pattern it was built from.
			"searchRegex", processor.SearchRegex.String(),
			"outputPattern", processor.OutputPattern,
		)
	}

	// Get the file names to process and compute the new names
	fileNames, err := getFolderFileNames(config.Folder)
	if err != nil {
		fatal(err.Error())
	}
	renumFolder := NewRenumFolder(config.SeasonNum, config.EpNum, fileNames, processors)
	for _, file := range renumFolder.RenumFiles {
		slog.Info("[Preview]", "oldName", file.OldName, "newName", file.NewName)
	}

	// Refuse a plan that would overwrite files, before touching anything
	plan := buildRenamePlan(renumFolder.RenumFiles)
	if err := validateRenamePlan(config.Folder, plan); err != nil {
		slog.Error("The following renames would destroy files:")
		for _, conflict := range joinedErrors(err) {
			slog.Error("  - " + conflict.Error())
		}
		fatal("Aborting, no file has been changed")
	}

	if config.DryRun {
		slog.Info("[DRY RUN] Exiting...")
		os.Exit(0)
	}

	if len(plan) == 0 {
		slog.Info("Nothing to rename")
		os.Exit(0)
	}

	// Ask for confirmation
	if !isOperationConfirmed(config.Force) {
		slog.Warn("Aborting the operation...")
		os.Exit(1)
	}

	// Rename files
	slog.Info("Continuing the operation...")
	for _, op := range plan {
		slog.Debug("[Rename]", "oldName", op.From, "newName", op.To)
	}
	if err := applyRenamePlan(config.Folder, plan); err != nil {
		for _, failure := range joinedErrors(err) {
			slog.Error(failure.Error())
		}
		fatal("Some files could not be renamed")
	}

	slog.Info(fmt.Sprintf("Renamed %d file(s)", len(plan)))
}

func getProcessors(config *Config) ([]*Processor, error) {
	if config.SearchPattern == "" {
		return getDefaultProcessors(), nil
	}

	processor, err := NewProcessor(config.SearchPattern, defaultOutputPattern)
	if err != nil {
		return nil, err
	}
	slog.Info("Using a custom search pattern", "searchPattern", config.SearchPattern)

	return []*Processor{processor}, nil
}

// joinedErrors takes apart what errors.Join gathered, so that each error is
// reported as its own message. Splitting the joined text on "\n" cut a message
// in two as soon as a file name held a newline — which a file name may.
func joinedErrors(err error) []error {
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		return joined.Unwrap()
	}

	return []error{err}
}

func isOperationConfirmed(force bool) bool {
	if force {
		slog.Info("Force mode enabled, continuing the operation...")
		return true
	}

	fmt.Print("Do you want to continue the operation? (y/N): ")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		slog.Debug("Unable to read the response, assuming it is 'n'", "error", err)
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y"
}

// getFolderFileNames lists the files of a folder, sub-folders excluded, in the
// order the episodes are numbered in.
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

	// os.ReadDir sorts lexicographically, which numbers "ep 10.mkv" before
	// "ep 2.mkv".
	sortNaturally(fileNames)

	return fileNames, nil
}
