package main

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// setupLogger installs the logger the flags ask for: text on the standard error
// output, JSON with --json, and the debug messages on top with --verbose.
func setupLogger(config *Config) {
	slog.SetDefault(slog.New(newLogHandler(config, os.Stderr)))
}

func newLogHandler(config *Config, output io.Writer) slog.Handler {
	options := &slog.HandlerOptions{
		Level:       logLevel(config.Verbose),
		ReplaceAttr: lowercaseLevel,
	}

	if config.Json {
		return slog.NewJSONHandler(output, options)
	}

	return slog.NewTextHandler(output, options)
}

func logLevel(verbose bool) slog.Level {
	if verbose {
		return slog.LevelDebug
	}

	return slog.LevelInfo
}

// lowercaseLevel spells the level the way logrus used to spell it. slog writes
// it in upper case, and whatever reads the --json output has been reading
// "info" since the first release.
func lowercaseLevel(groups []string, attr slog.Attr) slog.Attr {
	if len(groups) == 0 && attr.Key == slog.LevelKey {
		attr.Value = slog.StringValue(strings.ToLower(attr.Value.String()))
	}

	return attr
}

// fatal reports what went wrong and leaves with the exit code the README
// documents for a failure. slog has no fatal level, so the exit is ours to make
// — which also keeps it in one place instead of scattered os.Exit calls.
func fatal(message string, args ...any) {
	slog.Error(message, args...)
	os.Exit(1)
}
