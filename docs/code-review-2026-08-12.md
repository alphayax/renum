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

## Important

### 5. Companion files consumed episode numbers — `RenumFolder.go:19` — fixed

`Show S01E01.srt` became `Show S02E02.srt`, desynchronising every subtitle from
its video and shifting the rest of the folder.

**Fix:** files are grouped into episodes by their name without the extension and
share the episode number. A file that no processor matches is not an episode
either: it keeps its name and no longer consumes a number, so a stray
`cover.jpg` does not shift the folder. Only the last extension is dropped when
grouping, so `Show S01E01.fr.srt` still counts as its own episode — documented
in the README.

`RenumFolder.FolderPath`, the dead field listed below, was removed along the way.

### 6. The 4th default pattern ate the extension dot — `Processor.go:34` — fixed

`serie 1.mkv` became `serie_S01E01_mkv`, which no player will open: the
separators around the number were rewritten as underscores, the dot of the
extension included.

**Fix:** the separators are captured and put back as they were. The README
example is unaffected, `_1086_` still yields `_S12E01_`.

### 7. Every occurrence was substituted — `Processor.go:26` — fixed

`ReplaceAllString` rewrote each match of the pattern, so `S01E01 - rerun of
S01E01.mkv` became `S02E03 - rerun of S02E03.mkv`, and the separator pattern
renumbered the "2" of `serie 1 - part 2.mkv` as a second episode.

**Fix:** only the first match is replaced, through `FindStringSubmatchIndex` plus
`ExpandString` — which keeps the `$1` semantics the separator pattern needs.

### 8. The listing order numbered the wrong episodes — `main.go:113` — fixed

`os.ReadDir` sorts lexicographically, in which `ep 10.mkv` comes before
`ep 2.mkv`. Every folder whose numbers are not zero-padded was therefore
renumbered in the wrong order:

```
serie 1.mkv  serie 2.mkv  serie 3.mkv  serie 10.mkv  serie 11.mkv

serie 1.mkv  -> serie S04E01.mkv
serie 10.mkv -> serie S04E02.mkv     <- episode 10 numbered as the second
serie 11.mkv -> serie S04E03.mkv
serie 2.mkv  -> serie S04E04.mkv
```

**Fix:** the listing is sorted in the natural order (`NaturalOrder.go`), where
the digit runs of a name are compared by value. Everything that is not a number
keeps the byte order `os.ReadDir` used, and a run longer than an `int64` is
still compared correctly since the digits are compared as written, not parsed.

### Still open

Nothing left from the numbered list.

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
  `getFolderFileNames` — the last one being exactly where issue 1 lived. It is
  at 76% now, `main` itself being what is left untested.
- **Copy-paste bug in `docker.yaml:52` — fixed**: the Docker Hub description was
  pushed to `alphayax/chart-updater` instead of `alphayax/renum`, which failed
  the job at the end of every release. Once pointed at the right repository the
  step still failed with `Forbidden`, because `peter-evans/dockerhub-description`
  needs a Docker Hub token with the read/write/**delete** scope and
  `DOCKERHUB_TOKEN` did not have it. The scope was widened on 2026-08-25 and the
  whole Docker workflow is green on v1.1.1. **fixed**
- **Expired release token — fixed**: `RENUM_GITHUB_TOKEN` had expired, so
  GoReleaser built every binary and then failed to publish with a 401; no
  release had been created since v1.0.10 in January 2024. The job already
  declares `permissions: contents: write`, so it now uses the `GITHUB_TOKEN`
  GitHub injects into the run and there is no long-lived secret left to renew.
- **`.idea/` and `renum.iml` are versioned** while `.gitignore` only holds
  `dist/` — **fixed**: both are untracked and ignored, along with the binary
  `go build` drops in the repository root.
- **Go 1.20 is end of life** (`go.mod`, `Dockerfile`) — **fixed**: the module and
  the builder image are on 1.22.
- **logrus is in maintenance mode — fixed**: the logging moved to `log/slog`
  (`Logger.go`), which the standard library has carried since 1.21. The level is
  lower-cased to keep the `--json` output as it was, and `renum` now builds with
  no dependency at all.
- **Actions targeting Node 20 — fixed**: GitHub deprecated the runtime and every
  release ended with a warning. Both workflows are on the Node 24 majors
  (checkout v7, setup-go v7, goreleaser-action v7, docker/\*,
  dockerhub-description v5). goreleaser-action v7 pins GoReleaser to `~> v2`, so
  `.goreleaser.yaml` moved to `version: 2` and the deprecated `archives.format`
  became `formats`.
- `RenumFolder.FolderPath` was declared but never set — dead field, **removed**.
- `Config` relies on the global `flag` set, so calling `NewConfig()` twice in a
  single process panics, which makes the parsing hard to test properly. A
  dedicated `flag.NewFlagSet` would fix it.

---

## Suggested order

1. Critical issues 1, 2, 3, 4 — removes every data-loss path.
2. Test CI on pull requests, plus regression tests on the rename planner.
3. Issues 5 and 6 (subtitles, extension) — the most visible in daily use.
4. Cleanup: `filepath.Join`, `.gitignore`, Go 1.22, the Docker workflow.

## Releases

- **v1.1.0** (2026-08-12) — the critical fixes. Tagged before the companion-file
  fix landed, so `f3d2f86` shipped in no binary.
- **v1.1.1** (2026-08-25) — companion files and extensions, Go 1.22, the action
  bumps. First release where both workflows are green end to end: the GitHub
  release carries its eight archives and the Docker image, its tags and its
  Hub description are all published.
- **v1.1.2** (2026-08-25) — the natural order, the single substitution, the
  ignored IDE files.

## Found after the review

- **`renum -h` exits 1** and prints `invalid number of arguments` before the
  help: `Config.Parse` refuses an empty argument list before `main` gets to look
  at `config.Help`. The README documents `0` for a successful run, and asking for
  the help is not a failure. **open**
