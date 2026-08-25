package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// --json is a contract: whatever parses the output expects one JSON object per
// line, with the level spelled in lower case as logrus used to spell it.
func TestNewLogHandlerInJson(t *testing.T) {
	output := &bytes.Buffer{}
	logger := slog.New(newLogHandler(&Config{Json: true}, output))

	logger.Info("[Preview]", "oldName", "serie 1.mkv")

	var line map[string]any
	if err := json.Unmarshal(output.Bytes(), &line); err != nil {
		t.Fatalf("Expected a JSON line, but got %q (%v)", output.String(), err)
	}
	if line["level"] != "info" {
		t.Errorf("Expected the level to be 'info', but got %v", line["level"])
	}
	if line["msg"] != "[Preview]" {
		t.Errorf("Expected the message to be '[Preview]', but got %v", line["msg"])
	}
	if line["oldName"] != "serie 1.mkv" {
		t.Errorf("Expected oldName to be 'serie 1.mkv', but got %v", line["oldName"])
	}
}

func TestNewLogHandlerInText(t *testing.T) {
	output := &bytes.Buffer{}
	logger := slog.New(newLogHandler(&Config{}, output))

	logger.Warn("Aborting the operation...")

	line := output.String()
	if !strings.Contains(line, "level=warn") {
		t.Errorf("Expected the level to be 'warn', but got %q", line)
	}
	if !strings.Contains(line, `msg="Aborting the operation..."`) {
		t.Errorf("Expected the message in the line, but got %q", line)
	}
}

// --verbose is the only thing that opens the debug messages.
func TestNewLogHandlerLevels(t *testing.T) {
	cases := map[string]struct {
		config     *Config
		wantsDebug bool
	}{
		"default": {config: &Config{}, wantsDebug: false},
		"verbose": {config: &Config{Verbose: true}, wantsDebug: true},
	}

	for name, testCase := range cases {
		handler := newLogHandler(testCase.config, &bytes.Buffer{})
		if got := handler.Enabled(context.Background(), slog.LevelDebug); got != testCase.wantsDebug {
			t.Errorf("%s: expected the debug level enabled to be %t, but got %t",
				name, testCase.wantsDebug, got)
		}
		if !handler.Enabled(context.Background(), slog.LevelInfo) {
			t.Errorf("%s: expected the info level to be enabled", name)
		}
	}
}
