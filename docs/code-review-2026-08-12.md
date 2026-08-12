# Code review — 2026-08-12

Full review of the codebase at commit `32d2348` (main). Every issue in the
"Critical" section was reproduced by running the built binary; the reproduction
is quoted with each entry.

Status legend: **open** / **fixed** (with the commit or branch that fixed it).

---

## Critical

### 1. A sub-folder corrupts the whole file list — `main.go:118-129` — fixed

`getFolderFileNames` writes into `fileNames[i]` where `i` is the index in the
full directory listing (folders included), then truncates the slice to
`len(files) - folderCount`. As soon as a folder sorts *before* a file, every
following entry is shifted: one empty entry appears at the front, the episode
numbering is off by one, and the last file is dropped entirely.

```
Extras/  Show S01E01.mkv  Show S01E02.mkv  Show S01E03.mkv

[Preview] oldName=          newName=            <- empty entry
[Preview] oldName=S01E01 -> newName=S02E02      <- shifted numbering
[Preview] oldName=S01E02 -> newName=S02E03
                                                <- S01E03 silently dropped
```

On a real run this aborts with `fatal: rename .../t7/ .../t7/: file exists`
and exit code 1, after having possibly renamed part of the folder.

The demo GIF in the README does not show the bug only because its `Subs` folder
happens to sort *after* the media files.

**Fix:** build the slice with `append` instead of indexing by the listing
position.

### 2. Silent overwrite on rename collision — `main.go:85` — fixed

`os.Rename` silently replaces an existing target on Unix. Shifting episodes up
by one — an entirely ordinary operation — destroys a file:

```
before: Show S01E01.mkv = EPISODE-1
        Show S01E02.mkv = EPISODE-2

$ renum --season 1 --episode 2 --force .

after:  Show S01E03.mkv = EPISODE-1     <- EPISODE-2 is gone
```

**Fix:** a pre-flight validation pass that refuses to touch anything when two
files map to the same target, or when a target is an unrelated existing file;
plus a rename ordering that resolves chains and cycles inside the rename set
(via a temporary name) so that shifting a range up or down works losslessly.

### 3. `log.Fatal` in the middle of the rename loop — `main.go:89` — fixed

The first failing rename kills the process, leaving the folder half renamed,
with no rollback and no summary of what was and was not applied.

**Fix:** validate everything up front, then keep going on per-file errors,
report each one, and exit non-zero with a count at the end.

### 4. An invalid `--pattern` panics — `Processor.go:16` — fixed

`regexp.MustCompile` is applied to user input, so `renum --pattern "S[0-9" .`
prints a Go stack trace:

```
panic: regexp: Compile(`S[0-9`): error parsing regexp: missing closing ]: `[0-9`
```

**Fix:** `regexp.Compile` and a proper error message for user-supplied
patterns; keep the panicking constructor for the hard-coded defaults only.

---

## Important — open

| # | Issue | Location |
|---|---|---|
| 5 | **Companion files consume episode numbers.** `Show S01E01.srt` becomes `Show S02E02.srt`, desynchronising every subtitle from its video. Files should be grouped by stem, or filtered by extension. | `RenumFolder.go:19` |
| 6 | **The 4th default pattern eats the extension dot**: `serie 1.mkv` becomes `serie_S01E01_mkv`, which no player will open. | `Processor.go:34` |
| 7 | `ReplaceAllString` replaces *every* occurrence: `S01E01 - rerun S01E01.mkv` is substituted twice. A single replacement would be enough. | `Processor.go:26` |
| 8 | `os.ReadDir` sorts lexicographically, so `ep 10.mkv` comes before `ep 2.mkv`. Without natural sorting the numbering is wrong for any non zero-padded name. | `main.go:113` |
Issues 9 and 10 sat inside the code rewritten for the critical fixes, so they
were fixed along the way:

- **9 — fixed.** `fmt.Sprintf("%s/%s", …)` produced a double slash when the path
  ended with `/`, and the wrong separator on Windows — for which a binary is
  published. Now `filepath.Join`.
- **10 — fixed.** `os.Exit(-1)` was truncated to 255 on abort; it is now 1.

---

## Tooling and quality — open

- **No test CI.** Both workflows only trigger on tags; nothing runs on push or
  pull request. A `go test` + `go vet` + `golangci-lint` job on `pull_request`
  would have caught several of the issues above.
- **Coverage was 40%, with `main.go` at 0%** before this review: the three
  untested functions were `main`, `isOperationConfirmed` and
  `getFolderFileNames` — the last one being exactly where issue 1 lived.
- **Copy-paste bug in `docker.yaml:52`**: the Docker Hub description is pushed
  to `alphayax/chart-updater` instead of `alphayax/renum`.
- **`.idea/` and `renum.iml` are versioned** while `.gitignore` only holds
  `dist/`.
- **Go 1.20 is end of life** (`go.mod`, `Dockerfile`). Moving to 1.22+ is
  overdue; logrus is in maintenance mode and `log/slog` has been in the
  standard library since 1.21.
- `RenumFolder.FolderPath` is declared but never set — dead field.
- `Config` relies on the global `flag` set, so calling `NewConfig()` twice in a
  single process panics, which makes the parsing hard to test properly. A
  dedicated `flag.NewFlagSet` would fix it.

---

## Suggested order

1. Critical issues 1, 2, 3, 4 — removes every data-loss path.
2. Test CI on pull requests, plus regression tests on the rename planner.
3. Issues 5 and 6 (subtitles, extension) — the most visible in daily use.
4. Cleanup: `filepath.Join`, `.gitignore`, Go 1.22, the Docker workflow.
