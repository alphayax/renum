# Code review — 2026-08-25

Full review of the codebase at commit `52b7d0d` (main), the state left by the
[2026-08-12 review](./code-review-2026-08-12.md) and the releases that followed.
Every finding below was reproduced by running the built binary or by a probe
test; the reproduction is quoted with each entry.

Status legend: **open** / **fixed** (with the commit that fixed it).

---

## Where the project stands

The previous review left nothing open on its numbered list, and this one
confirms it: the data-loss paths are gone, and the checks that closed them hold
up under the folders that used to break them.

| Check | Result |
|---|---|
| `gofmt -l .` | clean |
| `go vet ./...` | clean |
| `golangci-lint run ./...` | 0 issues |
| `go test ./...` | pass, 79.9% of statements before this review, 82.9% after |

What was verified by hand, on real folders, all of it correct:

- a sub-folder in the listing (the old issue 1), including one sorting first;
- shifting a range up (`--episode 2` on `S01E01`–`S01E03`), the old issue 2 —
  every file kept;
- a target name held by a sub-folder — refused, folder untouched;
- a path that is a file, a missing folder, an empty folder — each reported with
  the exit code the README documents;
- `--dry-run`, `--json`, a trailing slash on the folder, `--episode 0`,
  `--season 100`;
- a folder mixing `serie 1..3` and `serie 10` with their subtitles and a
  `cover.jpg`: natural order, companion files on the same number, and the cover
  neither renamed nor counted.

The rest of this document is what is left.

---

## Findings

### 1. A file that cannot be renamed was left under a hidden temporary name — `RenamePlan.go:104-121` — fixed

`applyRenamePlan` parks a file under `.renum-tmp-0-<name>` to break a cycle. When
the rename it was parked for turned out to be impossible — because the file
holding the target failed to move away — the operation was reported as failed
and the parked file was left as it was. Nothing was lost, but the file was under
a dotted name nobody would look for, and the message named the temporary file
rather than the one the user knows:

```
error: unable to rename "ghost.mkv" to "z.mkv": no such file or directory
       unable to resolve the rename of ".renum-tmp-0-a.mkv" to "b.mkv"

folder: ".renum-tmp-0-a.mkv" holds "a.mkv"     <- a.mkv is gone from sight
        "b.mkv"              holds "b.mkv"
```

Reaching it takes a rename that fails half way through a batch — a permission, a
name too long for the file system once the prefix is added, a file pulled from
under Renum — since `validateRenamePlan` rules out the obvious conflicts before
anything is touched. It is the one path that survived the previous review's
"keep going on per-file errors" rule.

**Fix:** the parked file is put back under its own name when its rename is given
up, and the failure names that file. If something took its name in the meantime,
that is reported too, rather than silently left. `parked` was keyed by the target
of the rename; it is now `origins`, keyed by the temporary name, which is what
the restoration needs and what the guard was really asking about.

### 2. Error messages were cut on the newlines inside them — `main.go:64` and `main.go:92` — fixed

`errors.Join` glues its errors with a newline, and both call sites took the
joined text apart with `strings.Split(err.Error(), "\n")`. A file name may hold
a newline, so a single conflict could be logged as two half-messages — under
`--json`, two malformed lines for something that parses the output.

**Fix:** `joinedErrors` unwraps what `errors.Join` gathered instead of splitting
its text back.

### 3. `--help` wrote to the standard error output — `main.go:18` — fixed

`renum --help > help.txt` left the file empty, and `renum --help | less` showed
nothing: the usage went to the standard error output, the one the parsing errors
use. Asking for the help is not an error — the exit code has said so since the
previous review, the stream had not caught up.

**Fix:** `Config.SetOutput`, which main points at the standard output for the
help. A parsing error still reports itself, and prints the usage, on the standard
error output.

### 4. The confirmation prompt was untested — `main.go:115` — fixed

`isOperationConfirmed` is the last gate before any file is touched, and it was
one of the two functions at 0% coverage. It read `os.Stdin` directly, which is
what made it untestable.

**Fix:** it takes the reader to ask, main hands it `os.Stdin`, and the answers
are covered: `y`, `Y`, a padded `y`, `n`, an empty line, a closed input, and
anything else. All of them but a plain `y` decline — including an answer that
cannot be read, which is how Renum behaves in a container started without a
terminal.

---

## Open — by design, worth knowing

These are not defects to fix so much as consequences of what Renum is. They were
missing from the README; they are documented there now.

### 5. Anything carrying a number can be taken for an episode

The patterns only look at names, and the separator pattern `([_ ])[0-9]+([_ .])`
is a loose one. A poster is renumbered like an episode, and takes the number the
first episode was going to get:

```
$ renum --season 2 --dry-run .
[Preview] oldName="Season 1 Poster.jpg" newName="Season S02E01 Poster.jpg"
[Preview] oldName="Show S01E01.mkv"     newName="Show S02E02.mkv"
```

The previous review's fix — a file no processor matches is not an episode and
takes no number — covers `cover.jpg`, not `Season 1 Poster.jpg`. `--dry-run`
shows it before anything happens, and `--pattern` narrows the matching down;
tightening the default patterns would cost the loose names they exist for.

### 6. A symbolic link is renamed, not followed

A link is a file in the listing, so it is renumbered — and a link naming a file
of the same folder relatively ends up dangling, because that file has been
renamed too:

```
Show S01E01.mkv        -> Show S05E01.mkv
Show S01E02.mkv (link) -> Show S05E02.mkv -> Show S01E01.mkv   <- dangling
```

Rewriting link targets is a different tool. Documented as a limitation.

### 7. A name close to the file system's limit cannot be parked

Parking prepends 13 characters, so a 250-byte name cannot be parked on a file
system capped at 255. The rename is refused and reported, the file stays where
it is — and, since finding 1, under its own name. No fix beyond a shorter prefix
would buy much.

### 8. `main` itself is still untested

`main` is the last function at 0%: it is a sequence of `os.Exit` calls around
functions that are all covered. Testing it means running the binary, which the
review does by hand. A test that builds and drives it would close the gap and
pin the exit codes the README documents. Worth doing, not urgent.

---

## Tooling

Nothing found. For the record, what the previous review fixed is still standing:

- `test.yml` runs on every push to `main` and every pull request, with the Go
  version read from `go.mod`;
- the release and Docker workflows are on tags, with the built-in `GITHUB_TOKEN`
  and the right Docker Hub repository;
- `.idea/`, `renum.iml` and the binary `go build` drops in the root are ignored
  and untracked;
- Go 1.27 in `go.mod` and in the builder image, and still no dependency at all —
  `go.mod` is three lines and there is no `go.sum` to keep up to date.

Two cosmetic leftovers, neither of them worth a commit on its own:

- `.goreleaser.yaml` still carries the "This is an example" header and a
  `go generate ./...` hook the project has nothing to generate.
- The Dockerfile builds with `-a -installsuffix cgo`, both obsolete since Go
  1.10; `CGO_ENABLED=0` alone gives the same static binary and a faster build.

---

## Releases

The [2026-08-12 review](./code-review-2026-08-12.md) records the releases up to
v1.1.2. Since then:

- **v1.2.0** (2026-08-25) — the options on a flag set of renum's own, the
  logging moved from logrus to `log/slog`, and the test workflow running on
  every push and pull request.
- **v1.2.1** (2026-08-25) — Go 1.26 then 1.27, in `go.mod` and in the builder
  image. The state this review starts from.
- **v1.2.2** (2026-08-25) — the four fixes above, the limitations the README was
  missing, and the confirmation prompt under test.
