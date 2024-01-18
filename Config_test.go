package main

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.SeasonNum != 1 {
		t.Errorf("Expected SeasonNum to be 1, but got %d", config.SeasonNum)
	}

	if config.EpNum != 1 {
		t.Errorf("Expected EpNum to be 1, but got %d", config.EpNum)
	}

	if config.DryRun != false {
		t.Errorf("Expected DryRun to be false, but got %v", config.DryRun)
	}

	if config.Help != false {
		t.Errorf("Expected Help to be false, but got %v", config.Help)
	}

	if config.Force != false {
		t.Errorf("Expected Force to be false, but got %v", config.Force)
	}

	if config.Folder != "" {
		t.Errorf("Expected Folder to be empty, but got %s", config.Folder)
	}

}

func TestParseConfig(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config := NewConfig()

	args := []string{"-season=2", "--episode=3", "--dry-run", "--help", "--force", "testFolder"}
	os.Args = append(os.Args[:1], args...)

	if err := config.Parse(); err != nil {
		t.Fatal(err)
	}

	if config.SeasonNum != 2 {
		t.Errorf("Expected SeasonNum to be 2, but got %d", config.SeasonNum)
	}

	if config.EpNum != 3 {
		t.Errorf("Expected EpNum to be 3, but got %d", config.EpNum)
	}

	if config.DryRun != true {
		t.Errorf("Expected DryRun to be true, but got %v", config.DryRun)
	}

	if config.Help != true {
		t.Errorf("Expected Help to be true, but got %v", config.Help)
	}

	if config.Force != true {
		t.Errorf("Expected Force to be true, but got %v", config.Force)
	}

	if config.Folder != "testFolder" {
		t.Errorf("Expected Folder to be 'testFolder', but got %s", config.Folder)
	}
}

func TestParseConfigBad(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config := NewConfig()

	var args []string
	os.Args = append(os.Args[:1], args...)

	if err := config.Parse(); err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestUsage(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config := NewConfig()
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	config.Usage()
	str := buf.String()
	lines := strings.Split(str, "\n")
	if lines[0] != "Usage: renum [options] <folderPath>" {
		t.Errorf("Expected first line to be 'Usage: renum [options] <folderPath>', but got '%s'", lines[0])
	}
	if lines[1] != "Options:" {
		t.Errorf("Expected second line to be 'Options:', but got '%s'", lines[1])
	}
}
