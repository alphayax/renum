# Renum

Renum is a simple and efficient tool written in Go, designed to rename and renumber files in a directory. It's particularly useful for renaming series of files with a specific pattern.

## Features

- Rename files in a directory based on a specific pattern.
- Preview the changes before applying them.
- Easy to use with a simple command line interface.

## Installation

### Using pre-built Packages

`renum` is available for Windows, Linux and macOS. You can download the latest version from the [releases page](https://github.com/alphayax/renum/releases).

### Using Go

To install Renum, you need to have Go installed on your machine. Once you have Go installed, you can download and install Renum using the `go get` command:

```bash
go install github.com/alphayax/renum@latest
```

## Usage
To use Renum, navigate to the directory containing the files you want to rename and run the following command:
```bash
renum [options]
```
Here are the available options: 
- `--season`: The season number to start from.
- `--episode`: The episode number to start from.
- `--dry-run`: Preview the changes without applying them.

## Examples

Let's say you have a directory containing the following files:
```
XXX-Fansub]_Xxx_Xxxxx_1086_[VOSTFR][FHD_1920x1080].xxx
XXX-Fansub]_Xxx_Xxxxx_1087_[VOSTFR][FHD_1920x1080].xxx
XXX-Fansub]_Xxx_Xxxxx_1088_[VOSTFR][FHD_1920x1080].xxx
```

To rename these files, you can run the following command:
```bash
renum --season 12 --episode 1 /path/to/directory
```

This will rename the files to:
```
[XXX-Fansub]_Xxx_Xxxxx_S12E01_[VOSTFR][FHD_1920x1080].xxx
[XXX-Fansub]_Xxx_Xxxxx_S12E02_[VOSTFR][FHD_1920x1080].xxx
[XXX-Fansub]_Xxx_Xxxxx_S12E03_[VOSTFR][FHD_1920x1080].xxx
```


## Testing
To run the tests for Renum, navigate to the project directory and run the following command:
```bash
go test ./...
```


## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

