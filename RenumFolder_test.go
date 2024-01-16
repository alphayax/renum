package main

import (
	"testing"
)

func TestNewRenumFolder(t *testing.T) {
	seasonNum := uint(2)
	epNum := uint(3)
	folder := "test_folder"
	fileNames := []string{"S01E01.mkv", "S01E02.mkv"}

	renumFolder := NewRenumFolder(seasonNum, epNum, folder, fileNames)

	if renumFolder.SeasonNum != seasonNum {
		t.Errorf("Expected SeasonNum to be %d, but got %d", seasonNum, renumFolder.SeasonNum)
	}

	if renumFolder.EpNum != epNum {
		t.Errorf("Expected EpNum to be %d, but got %d", epNum, renumFolder.EpNum)
	}

	if renumFolder.FolderPath != folder {
		t.Errorf("Expected FolderPath to be %s, but got %s", folder, renumFolder.FolderPath)
	}

	if len(renumFolder.RenumFiles) != len(fileNames) {
		t.Errorf("Expected RenumFiles length to be %d, but got %d", len(fileNames), len(renumFolder.RenumFiles))
	}

	if renumFolder.RenumFiles[0].NewName != "S02E03.mkv" {
		t.Errorf("Expected RenumFiles[0].NewName to be S02E03.mkv, but got %s", renumFolder.RenumFiles[0].NewName)
	}
	if renumFolder.RenumFiles[1].NewName != "S02E04.mkv" {
		t.Errorf("Expected RenumFiles[0].NewName to be S02E04.mkv, but got %s", renumFolder.RenumFiles[0].NewName)
	}
}

/*
func TestPreview(t *testing.T) {
	// This function prints output to the console, so it's hard to test without changing the function to return a value.
	// You could potentially redirect the output and check that, but it's generally not recommended to test log outputs.
}

func TestRename(t *testing.T) {
	// This function interacts with the file system, so it's hard to test without creating actual files and folders.
	// You could potentially create a temporary directory and files, run the function, and then check the results.
	// Remember to clean up any temporary files and directories after the test.

	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	seasonNum := uint(2)
	epNum := uint(3)
	folder := tmpDir
	fileNames := []string{"S01E01.mkv", "S01E02.mkv"}

	// Create test files
	for _, fileName := range fileNames {
		filePath := filepath.Join(tmpDir, fileName)
		if err := ioutil.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	renumFolder := NewRenumFolder(seasonNum, epNum, folder, fileNames)
	renumFolder.Rename()

	// Check that the files have been renamed
	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	for i, file := range files {
		expectedName := renumFolder.RenumFiles[i].NewName
		if file.Name() != expectedName {
			t.Errorf("Expected file name to be %s, but got %s", expectedName, file.Name())
		}
	}
}
*/
