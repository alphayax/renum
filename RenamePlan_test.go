package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// makeFolder creates a temporary folder holding one file per entry, the file
// content being its name, so that a lost or overwritten file can be detected.
func makeFolder(t *testing.T, names ...string) string {
	t.Helper()

	folder := t.TempDir()
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(folder, name), []byte(name), 0o600); err != nil {
			t.Fatalf("unable to create the fixture %q: %v", name, err)
		}
	}

	return folder
}

// folderContent maps each file name to its original content.
func folderContent(t *testing.T, folder string) map[string]string {
	t.Helper()

	entries, err := os.ReadDir(folder)
	if err != nil {
		t.Fatalf("unable to read %q: %v", folder, err)
	}

	content := make(map[string]string, len(entries))
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(folder, entry.Name()))
		if err != nil {
			t.Fatalf("unable to read %q: %v", entry.Name(), err)
		}
		content[entry.Name()] = string(data)
	}

	return content
}

func TestBuildRenamePlanSkipsUnchangedNames(t *testing.T) {
	files := []*RenumFile{
		{OldName: "a.mkv", NewName: "a.mkv"},
		{OldName: "b.mkv", NewName: "c.mkv"},
	}

	plan := buildRenamePlan(files)
	if len(plan) != 1 {
		t.Fatalf("Expected 1 rename, but got %d", len(plan))
	}
	if plan[0].From != "b.mkv" || plan[0].To != "c.mkv" {
		t.Errorf("Expected b.mkv -> c.mkv, but got %s -> %s", plan[0].From, plan[0].To)
	}
}

func TestValidateRenamePlanAcceptsAShift(t *testing.T) {
	folder := makeFolder(t, "S01E01.mkv", "S01E02.mkv")
	plan := []renameOp{
		{From: "S01E01.mkv", To: "S01E02.mkv"},
		{From: "S01E02.mkv", To: "S01E03.mkv"},
	}

	if err := validateRenamePlan(folder, plan); err != nil {
		t.Errorf("Expected the shift to be valid, but got %v", err)
	}
}

func TestValidateRenamePlanRejectsDuplicateTargets(t *testing.T) {
	folder := makeFolder(t, "a.mkv", "b.mkv")
	plan := []renameOp{
		{From: "a.mkv", To: "S01E01.mkv"},
		{From: "b.mkv", To: "S01E01.mkv"},
	}

	if err := validateRenamePlan(folder, plan); err == nil {
		t.Errorf("Expected two files claiming the same name to be rejected")
	}
}

func TestValidateRenamePlanRejectsOverwritingAnUntouchedFile(t *testing.T) {
	folder := makeFolder(t, "a.mkv", "keep.mkv")
	plan := []renameOp{
		{From: "a.mkv", To: "keep.mkv"},
	}

	if err := validateRenamePlan(folder, plan); err == nil {
		t.Errorf("Expected overwriting an untouched file to be rejected")
	}
}

func TestValidateRenamePlanRejectsNamesEscapingTheFolder(t *testing.T) {
	folder := makeFolder(t, "a.mkv")
	plan := []renameOp{
		{From: "a.mkv", To: "../a.mkv"},
	}

	if err := validateRenamePlan(folder, plan); err == nil {
		t.Errorf("Expected a name escaping the folder to be rejected")
	}
}

// Shifting episodes up used to overwrite the next file: every rename targeted a
// name still held by the file processed after it.
func TestApplyRenamePlanShiftUpKeepsEveryFile(t *testing.T) {
	folder := makeFolder(t, "S01E01.mkv", "S01E02.mkv", "S01E03.mkv")
	plan := []renameOp{
		{From: "S01E01.mkv", To: "S01E02.mkv"},
		{From: "S01E02.mkv", To: "S01E03.mkv"},
		{From: "S01E03.mkv", To: "S01E04.mkv"},
	}

	if err := validateRenamePlan(folder, plan); err != nil {
		t.Fatalf("Expected the plan to be valid, but got %v", err)
	}
	if err := applyRenamePlan(folder, plan); err != nil {
		t.Fatalf("Expected the plan to apply, but got %v", err)
	}

	expected := map[string]string{
		"S01E02.mkv": "S01E01.mkv",
		"S01E03.mkv": "S01E02.mkv",
		"S01E04.mkv": "S01E03.mkv",
	}
	content := folderContent(t, folder)
	if len(content) != len(expected) {
		t.Fatalf("Expected %d files, but got %d: %v", len(expected), len(content), content)
	}
	for name, want := range expected {
		if content[name] != want {
			t.Errorf("Expected %q to hold %q, but got %q", name, want, content[name])
		}
	}
}

func TestApplyRenamePlanShiftDownKeepsEveryFile(t *testing.T) {
	folder := makeFolder(t, "S01E02.mkv", "S01E03.mkv")
	plan := []renameOp{
		{From: "S01E02.mkv", To: "S01E01.mkv"},
		{From: "S01E03.mkv", To: "S01E02.mkv"},
	}

	if err := applyRenamePlan(folder, plan); err != nil {
		t.Fatalf("Expected the plan to apply, but got %v", err)
	}

	content := folderContent(t, folder)
	if content["S01E01.mkv"] != "S01E02.mkv" || content["S01E02.mkv"] != "S01E03.mkv" {
		t.Errorf("Unexpected folder content: %v", content)
	}
}

// A swap is a cycle: no rename can start until one file is parked aside.
func TestApplyRenamePlanResolvesACycle(t *testing.T) {
	folder := makeFolder(t, "a.mkv", "b.mkv")
	plan := []renameOp{
		{From: "a.mkv", To: "b.mkv"},
		{From: "b.mkv", To: "a.mkv"},
	}

	if err := validateRenamePlan(folder, plan); err != nil {
		t.Fatalf("Expected the swap to be valid, but got %v", err)
	}
	if err := applyRenamePlan(folder, plan); err != nil {
		t.Fatalf("Expected the swap to apply, but got %v", err)
	}

	content := folderContent(t, folder)
	if len(content) != 2 {
		t.Fatalf("Expected 2 files, but got %v", content)
	}
	if content["a.mkv"] != "b.mkv" || content["b.mkv"] != "a.mkv" {
		t.Errorf("Expected the files to be swapped, but got %v", content)
	}
}

func TestApplyRenamePlanReportsAMissingSource(t *testing.T) {
	folder := makeFolder(t, "a.mkv")
	plan := []renameOp{
		{From: "ghost.mkv", To: "S01E01.mkv"},
		{From: "a.mkv", To: "S01E02.mkv"},
	}

	err := applyRenamePlan(folder, plan)
	if err == nil {
		t.Fatalf("Expected the missing source to be reported")
	}

	// The failure must not stop the renames that can still be applied.
	content := folderContent(t, folder)
	if content["S01E02.mkv"] != "a.mkv" {
		t.Errorf("Expected a.mkv to be renamed despite the failure, but got %v", content)
	}
}

// A rename that can never be applied used to leave its file parked under the
// hidden temporary name it had been moved to: the file was still there, but
// under a name nobody would look for. Here "b.mkv" is never moved away — the
// rename that would have freed it fails — so "a.mkv" has to stay "a.mkv".
func TestApplyRenamePlanPutsAnUnresolvableFileBack(t *testing.T) {
	folder := makeFolder(t, "a.mkv", "b.mkv")
	plan := []renameOp{
		{From: "a.mkv", To: "b.mkv"},
		{From: "ghost.mkv", To: "z.mkv"},
	}

	if err := applyRenamePlan(folder, plan); err == nil {
		t.Fatalf("Expected the impossible rename to be reported")
	}

	content := folderContent(t, folder)
	if len(content) != 2 {
		t.Fatalf("Expected 2 files, but got %v", content)
	}
	if content["a.mkv"] != "a.mkv" || content["b.mkv"] != "b.mkv" {
		t.Errorf("Expected both files under their own name, but got %v", content)
	}
}

// The failure must name the file the user knows, not the temporary name it was
// parked under along the way.
func TestApplyRenamePlanReportsTheOriginalName(t *testing.T) {
	folder := makeFolder(t, "a.mkv", "b.mkv")
	plan := []renameOp{
		{From: "a.mkv", To: "b.mkv"},
		{From: "ghost.mkv", To: "z.mkv"},
	}

	err := applyRenamePlan(folder, plan)
	if err == nil {
		t.Fatalf("Expected an error")
	}
	if !strings.Contains(err.Error(), `unable to resolve the rename of "a.mkv" to "b.mkv"`) {
		t.Errorf("Expected the original name in the failure, but got %v", err)
	}
	if strings.Contains(err.Error(), tempPrefix) {
		t.Errorf("Expected no temporary name in the failure, but got %v", err)
	}
}

func TestApplyRenamePlanLeavesNoTemporaryFile(t *testing.T) {
	folder := makeFolder(t, "a.mkv", "b.mkv", "c.mkv")
	plan := []renameOp{
		{From: "a.mkv", To: "b.mkv"},
		{From: "b.mkv", To: "c.mkv"},
		{From: "c.mkv", To: "a.mkv"},
	}

	if err := applyRenamePlan(folder, plan); err != nil {
		t.Fatalf("Expected the rotation to apply, but got %v", err)
	}

	var names []string
	for name := range folderContent(t, folder) {
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) != 3 || names[0] != "a.mkv" || names[1] != "b.mkv" || names[2] != "c.mkv" {
		t.Errorf("Expected exactly a/b/c.mkv, but got %v", names)
	}
}
