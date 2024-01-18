# Renum

[![GitHub release](https://img.shields.io/github/v/release/alphayax/renum)](https://github.com/alphayax/renum/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/alphayax/renum)](https://goreportcard.com/report/github.com/alphayax/renum)

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

### Using Docker

You can also use Renum with Docker. To do so, you can run the following command:

```bash
docker run --rm -it -v /path/to/directory:/data alphayax/renum:latest [options] /data
```

## Usage
To use Renum, run the following command by passing the path to the directory containing the files you want to rename as last argument:
```bash
renum [options] /path/to/directory
```
### Options
- `-s <NUM>`, `--season <NUM>`: The season number to use.
- `-e <NUM>`, `--episode <NUM>`: The episode number to start from. Will be incremented for each file.
- `-h`, `--help`: Display the help message.
- `--force`: Don't ask for confirmation before applying the changes.
- `--dry-run`: Preview the changes without applying them.

### Filename patterns detected
- `S[0-9]+E[0-9]+`: containing `S1E01` or `S01E01`.
- ` [0-9]{1,2}x[0-9]+ `: containing ` 1x01 ` or ` 01x01 `.
- `^E[0-9]+`: starting by `E01` or `E001`...
- `[_ ][0-9]+[_ .]`: containing `_01_` or `_001_` or `_0001_` or ` 01 ` or `001`...

## Examples

Let's say you have a directory containing the following files:
```
[XXX-Fansub]_Xxx_Xxxxx_1086_[VOSTFR][FHD_1920x1080].xxx
[XXX-Fansub]_Xxx_Xxxxx_1087_[VOSTFR][FHD_1920x1080].xxx
[XXX-Fansub]_Xxx_Xxxxx_1088_[VOSTFR][FHD_1920x1080].xxx
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


## Sponsoring
Feel free to send crypto donations to the following addresses:
- Solana (SOL): `HUC9MmKR6iCtxu25h8hsgnVqXzeQMTK9ThQSLMFYNJBC`
- Ethereum (ETH): `0xc12Ef701Dd7e5060f441b30fE569D8D7E8a230a7`
- Bitcoin (BTC): `bc1qv7g3d8u9svn4w0pzfjafa7jzyglwjfkzjuc73g`
