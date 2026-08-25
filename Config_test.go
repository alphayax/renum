package main

import (
	"bytes"
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

// The options used to be declared on the global flag set, which panics on the
// second declaration of the same option in a process — the reason every test
// here had to rebuild flag.CommandLine and rewrite os.Args before touching a
// Config.
func TestNewConfigTwiceInTheSameProcess(t *testing.T) {
	first := NewConfig()
	second := NewConfig()

	if err := first.Parse([]string{"--season=2", "first"}); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if err := second.Parse([]string{"--season=3", "second"}); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if first.SeasonNum != 2 || first.Folder != "first" {
		t.Errorf("Expected the first config to keep its own values, but got %d and %q",
			first.SeasonNum, first.Folder)
	}
	if second.SeasonNum != 3 || second.Folder != "second" {
		t.Errorf("Expected the second config to keep its own values, but got %d and %q",
			second.SeasonNum, second.Folder)
	}
}

func TestParseConfig(t *testing.T) {
	config := NewConfig()

	args := []string{"-season=2", "--episode=3", "--dry-run", "--help", "--force", "testFolder"}
	if err := config.Parse(args); err != nil {
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
	config := NewConfig()

	if err := config.Parse(nil); err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

// "renum -h" used to be answered with "invalid number of arguments" and exit
// code 1: the missing folder was checked before the help was even looked at.
func TestParseConfigHelpWithoutAFolder(t *testing.T) {
	for _, arg := range []string{"-h", "--help"} {
		config := NewConfig()

		if err := config.Parse([]string{arg}); err != nil {
			t.Errorf("Expected %q to be accepted, but got %v", arg, err)
		}
		if !config.Help {
			t.Errorf("Expected %q to ask for the help", arg)
		}
	}
}

// An unknown option is an error of ours to report: the flag set must not print
// anything of its own on the way out.
func TestParseConfigUnknownOption(t *testing.T) {
	output := &bytes.Buffer{}
	config := newConfig(output)

	if err := config.Parse([]string{"--nope", "testFolder"}); err == nil {
		t.Errorf("Expected an error on an unknown option, but got nil")
	}
	if output.Len() != 0 {
		t.Errorf("Expected the flag set to stay silent, but it wrote %q", output.String())
	}
}

// "renum --help > help.txt" used to write the help on the standard error output
// and leave the file empty.
func TestSetOutput(t *testing.T) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	config := newConfig(stderr)

	config.SetOutput(stdout)
	config.Usage()

	if stdout.Len() == 0 {
		t.Errorf("Expected the usage on the chosen output, but it was empty")
	}
	if stderr.Len() != 0 {
		t.Errorf("Expected nothing on the first output, but it wrote %q", stderr.String())
	}
}

func TestUsage(t *testing.T) {
	output := &bytes.Buffer{}
	config := newConfig(output)

	config.Usage()

	lines := strings.Split(output.String(), "\n")
	if lines[0] != "Usage: renum [options] <folderPath>" {
		t.Errorf("Expected first line to be 'Usage: renum [options] <folderPath>', but got '%s'", lines[0])
	}
	if lines[1] != "Options:" {
		t.Errorf("Expected second line to be 'Options:', but got '%s'", lines[1])
	}
	if !strings.Contains(output.String(), "-pattern") {
		t.Errorf("Expected the options to be listed, but got '%s'", output.String())
	}
}
