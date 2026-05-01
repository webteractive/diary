# Diary Go Implementation Plan

Date: 2026-05-01
Status: Draft
Source PRD: `docs/plans/2026-05-01-ai-harness-journal-cli-prd.md`

## Objective

Implement the first usable version of `diary`, a local-first Go CLI for recording and retrieving compact AI harness context.

V1 commands:

- `diary record`
- `diary get`
- `diary list`
- `diary install-skills`
- `diary self-update`

The tool should store project-aware Markdown records under `.diary/`, maintain rebuildable JSON indexes, and support exact retrieval by id or hash.

## Implementation Principles

- Prefer Go standard library where practical.
- Use Cobra for CLI command structure.
- Use `gopkg.in/yaml.v3` for config/frontmatter parsing unless a simpler manual renderer proves sufficient.
- Keep Markdown files as the source of truth.
- Treat `index.json` as a cache that can be rebuilt.
- Keep command output predictable for AI harnesses.
- Keep write operations simple and append-friendly.
- Do not read `.env` files.

## Proposed Project Layout

```text
cmd/diary/
  main.go
internal/cli/
  root.go
  record.go
  get.go
  list.go
  install_skills.go
internal/install/
  install.go
  templates.go
internal/update/
  update.go
internal/project/
  resolve.go
internal/storage/
  paths.go
  record.go
  index.go
  markdown.go
internal/hash/
  hash.go
internal/render/
  markdown.go
  json.go
```

## Phase 1: Scaffold Go Module

Tasks:

- Add `go.mod`.
- Add dependencies: `github.com/spf13/cobra` and likely `gopkg.in/yaml.v3`.
- Add `cmd/diary/main.go`.
- Add internal package directories.
- Add Cobra root command plus `record`, `get`, `list`, `install-skills`, and `self-update` subcommands.

Acceptance:

- `go test ./...` passes.
- `go run ./cmd/diary --help` prints available commands.
- Unknown commands return a non-zero exit code and concise error.
- Shared flags such as `--project` and `--json` are available where relevant.

## Phase 2: Project Resolution

Implement project resolution order:

1. `--project <name>`.
2. `.diary/config.yml` project name, if present.
3. Git repository root directory basename, if inside Git.
4. Current working directory basename.

Tasks:

- Add `internal/project.Resolve`.
- Find Git root by walking parent directories for `.git`.
- Sanitize project names for safe directory use.

Acceptance:

- Tests cover explicit project, config project, Git root fallback, and cwd fallback.
- Project names are stable when running from subdirectories.

## Phase 3: Record Command

Command examples:

```bash
diary record "Implemented storage layout"
diary record --project campaign-builder "Updated importer validation flow" --file app/Import.php
echo "Next run should inspect storage/index.go" | diary record
```

Tasks:

- Parse message from args or stdin.
- Create `.diary/` automatically.
- Create `.diary/projects/<project>/records/`.
- Generate ids using UTC timestamp plus short random suffix.
- Write Markdown files with YAML frontmatter.
- Compute and store `sha256:<hash>`.
- Update `latest.md`.
- Update or rebuild project `index.json`.

Acceptance:

- Empty messages fail clearly.
- Recording creates the expected directory structure.
- Markdown can be inspected without the CLI.
- `latest.md` points to the newest useful context.
- `index.json` includes id, hash, timestamp, project, preview, files, refs, and tags.

## Phase 4: Hashing and Normalization

Tasks:

- Define canonical hash input.
- Exclude the `hash` field from hash input.
- Normalize line endings to `\n`.
- Sort frontmatter keys before hashing.

Acceptance:

- Hashes are deterministic across repeated reads.
- Changing record body changes the hash.
- Rewriting the same logical record does not change the hash because of map ordering.

## Phase 5: Get Command

Command examples:

```bash
diary get
diary get --project campaign-builder
diary get --id 2026-05-01T103000Z-codex-a7f3c9
diary get --hash abc123
diary get --json
```

Tasks:

- Resolve project.
- Return exact record for `--id`.
- Return exact record for unambiguous `--hash` prefix.
- Return latest compact context by default.
- Respect an output budget flag, for example `--max-chars`.
- Support Markdown default output and JSON output.

Acceptance:

- Default `get` returns prompt-ready context.
- Exact id lookup returns one record.
- Ambiguous hash prefixes fail with matching candidates.
- Missing records produce actionable errors.

## Phase 6: List Command

Command examples:

```bash
diary list
diary list --projects
diary list --project campaign-builder --json
```

Tasks:

- List records for the resolved project by default.
- List known projects with `--projects`.
- Show id, hash prefix, timestamp, project, preview, status, and tags.
- Support JSON output.

Acceptance:

- `list` does not emit full context.
- `list --projects` works without resolving a project first.
- Output is stable enough for harness parsing.

## Phase 7: Tests and Quality Pass

## Phase 7: Install Skills Command

Command examples:

```bash
diary install-skills --target codex
diary install-skills --target claude
diary install-skills --target all --dry-run
```

Tasks:

- Add embedded Diary skill templates for `diary-get`, `diary-record`, and `diary-list`.
- Install to default harness skill directories.
- Support `--path`, `--force`, `--dry-run`, and `--json`.
- Refuse to overwrite existing `SKILL.md` without `--force`.

Acceptance:

- Dry runs do not write files.
- Existing files are preserved unless `--force` is supplied.
- Custom paths work for single targets.
- `--target all` installs both supported targets.
- Each target receives `diary-get/SKILL.md`, `diary-record/SKILL.md`, and `diary-list/SKILL.md`.

## Phase 8: Tests and Quality Pass

## Phase 8: Self Update Command

Command examples:

```bash
diary self-update
diary self-update --version v0.0.1
diary self-update --dry-run
```

Tasks:

- Resolve the latest GitHub release unless a version is supplied.
- Detect OS and architecture.
- Download the matching release archive.
- Replace the running binary on macOS and Linux.
- Return a clear unsupported message on Windows.
- Support `--dry-run` and `--json`.

Acceptance:

- Dry runs do not download or replace files.
- Version-specific URLs match release asset naming.
- Latest release resolution is tested.
- Windows returns a clear unsupported error.

## Phase 9: Tests and Quality Pass

Add tests for:

- CLI parsing.
- Project resolution.
- Storage path generation.
- Markdown/frontmatter rendering.
- Hash determinism.
- Record creation.
- Index update/rebuild.
- Get default, id, and hash retrieval.
- List records and projects.
- Install-skills dry-run, force, existing-file protection, and custom path handling.
- Self-update dry-run, version URL construction, latest release parsing, and Windows unsupported behavior.
- Error messages for missing journal, empty message, and ambiguous hash.

Run before handoff:

```bash
go test ./...
go run ./cmd/diary record "Smoke test"
go run ./cmd/diary list
go run ./cmd/diary get
```

## Open Implementation Decisions

- Module path for `go.mod`.
- Exact YAML/frontmatter rendering strategy with `yaml.v3`.
- Whether `latest.md` should duplicate the latest record body or contain a pointer plus rendered summary.
- Whether `index.md` is still needed once project-level `index.json` exists.

## Recommended First Build Slice

Start with the smallest end-to-end path:

1. `go mod init`.
2. Implement `diary record "message"`.
3. Write one Markdown record under `.diary/projects/<project>/records/`.
4. Implement `diary list`.
5. Implement `diary get` returning latest.
6. Add hash/id lookup.

This gives a working tool quickly while keeping later features grounded in real files.
