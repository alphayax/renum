package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// tempPrefix is used to park a file whose target name is still held by another
// file of the same batch.
const tempPrefix = ".renum-tmp-"

// renameOp is a single rename to apply inside the target folder. From and To are
// bare file names, never paths.
type renameOp struct {
	From string
	To   string
}

// buildRenamePlan lists the renames to apply. Files whose name is left unchanged
// are skipped: renaming them onto themselves is pure noise.
func buildRenamePlan(files []*RenumFile) []renameOp {
	plan := make([]renameOp, 0, len(files))
	for _, file := range files {
		if file.OldName == file.NewName {
			continue
		}
		plan = append(plan, renameOp{From: file.OldName, To: file.NewName})
	}

	return plan
}

// validateRenamePlan reports every conflict that would make the plan lose data.
// It must be called before a single file is touched: os.Rename silently replaces
// its target, so an unchecked plan can destroy files without any warning.
func validateRenamePlan(folderPath string, plan []renameOp) error {
	sources := make(map[string]bool, len(plan))
	for _, op := range plan {
		sources[op.From] = true
	}

	var conflicts []error
	targets := make(map[string]string, len(plan))
	for _, op := range plan {
		// A pattern must never produce a name that escapes the folder.
		if filepath.Base(op.To) != op.To || filepath.Base(op.From) != op.From {
			conflicts = append(conflicts, fmt.Errorf("%q would be renamed to the invalid name %q", op.From, op.To))
			continue
		}

		// Two files claiming the same name: one would overwrite the other.
		if previous, taken := targets[op.To]; taken {
			conflicts = append(conflicts, fmt.Errorf("%q and %q would both be renamed to %q", previous, op.From, op.To))
			continue
		}
		targets[op.To] = op.From

		// An existing file that the batch never moves away would be overwritten.
		if sources[op.To] {
			continue
		}
		if _, err := os.Lstat(filepath.Join(folderPath, op.To)); err == nil {
			conflicts = append(conflicts, fmt.Errorf("renaming %q to %q would overwrite an existing file", op.From, op.To))
		}
	}

	return errors.Join(conflicts...)
}

// applyRenamePlan performs the renames. A target may still be held by another
// file of the batch (shifting a range of episodes up does exactly that), so
// blocked renames are retried, and a cycle is broken by parking one file under a
// temporary name. Individual failures are collected instead of aborting the run
// half way through.
func applyRenamePlan(folderPath string, plan []renameOp) error {
	pending := make([]renameOp, len(plan))
	copy(pending, plan)

	parked := make(map[string]bool, len(plan))
	var failures []error

	for len(pending) > 0 {
		var blocked []renameOp
		progress := false

		for _, op := range pending {
			if _, err := os.Lstat(filepath.Join(folderPath, op.To)); err == nil {
				// The target is still held by a file we have not moved yet.
				blocked = append(blocked, op)
				continue
			}
			if err := renameOne(folderPath, op); err != nil {
				failures = append(failures, err)
			}
			progress = true
		}

		if len(blocked) == 0 {
			break
		}

		if !progress {
			// Every remaining rename waits for another one: the batch contains a
			// cycle. Parking one file frees its name and unblocks the chain.
			op := blocked[0]
			if parked[op.To] {
				failures = append(failures, fmt.Errorf("unable to resolve the rename of %q to %q", op.From, op.To))
				blocked = blocked[1:]
			} else {
				tempName, err := parkFile(folderPath, op.From)
				if err != nil {
					failures = append(failures, err)
					blocked = blocked[1:]
				} else {
					parked[op.To] = true
					blocked[0].From = tempName
				}
			}
		}

		pending = blocked
	}

	return errors.Join(failures...)
}

func renameOne(folderPath string, op renameOp) error {
	from := filepath.Join(folderPath, op.From)
	to := filepath.Join(folderPath, op.To)
	if err := os.Rename(from, to); err != nil {
		return fmt.Errorf("unable to rename %q to %q: %w", op.From, op.To, err)
	}

	return nil
}

// parkFile moves a file to a free temporary name inside the same folder.
func parkFile(folderPath string, name string) (string, error) {
	for i := 0; i < 1000; i++ {
		tempName := fmt.Sprintf("%s%d-%s", tempPrefix, i, name)
		if _, err := os.Lstat(filepath.Join(folderPath, tempName)); err == nil {
			continue
		}
		if err := os.Rename(filepath.Join(folderPath, name), filepath.Join(folderPath, tempName)); err != nil {
			return "", fmt.Errorf("unable to move %q out of the way: %w", name, err)
		}

		return tempName, nil
	}

	return "", fmt.Errorf("unable to find a free temporary name for %q", name)
}
