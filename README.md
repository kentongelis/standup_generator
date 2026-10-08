# standup_generator

A CLI tool that generates a daily standup update from your local git commits and GitHub activity. See `proposal.md` for the full design.

## Status

MVP complete:

- **Git collector** — reads commits authored by you from a list of local repos (all branches) via `go-git`
- **GitHub collector** — uses the GitHub Search API to find PRs you opened or merged, PRs you reviewed (excluding your own), and assigned issues that were closed
- **Concurrent collection** — all collectors run in parallel; if one fails, the others still report and a warning naming the failed source is printed to stderr
- **Last-workday logic** — activity is collected since midnight of the previous weekday (on Monday, that's Friday)
- **Template output** — activity is split into *Yesterday* / *Today* sections and grouped by repo
- **`--copy`** — copies the generated standup to your clipboard
- **Config file** — `~/standup.yaml`, with a `GITHUB_TOKEN` env override

Planned next (see `proposal.md`): AI-written summaries (`--ai`), Slack posting, and blocker detection.

## Installation

Requires Go 1.27+.

```bash
go install github.com/kentongelis/standup@latest
# or, from a clone:
go build -o standup .
```

## Configuration

Create `~/standup.yaml`:

```yaml
github_username: your-github-username
github_token: ghp_...          # or set the GITHUB_TOKEN env var instead
git_email: you@example.com     # optional; defaults to `git config --global user.email`
repos:
  - /path/to/repo-one
  - /path/to/repo-two
```

| Key | Purpose |
|---|---|
| `github_username` | GitHub account to search activity for (required for the GitHub collector) |
| `github_token` | Personal access token; `GITHUB_TOKEN` in the environment overrides it |
| `git_email` | Author email used to filter local commits |
| `repos` | Absolute paths to local git repos to scan |

The config file is optional — without it, the git collector falls back to your global git email and the GitHub collector reports a warning that `github_username` is not set.

Because the file can hold your token, restrict its permissions:

```bash
chmod 600 ~/standup.yaml
```

## Usage

```bash
standup          # print your standup
standup --copy   # print it and copy it to the clipboard
```

Example output:

```
Yesterday
 standup
     • Committed: added copy to clipboard flag
     • Merged PR #12 Add GitHub collector

Today
 • Nothing recorded
```

If a source fails (for example, no network), you'll see a warning like `warning (github): ...` on stderr and the standup is still generated from the remaining sources.

## Development

```bash
go build ./...   # build
go run .         # run the CLI
go vet ./...     # static checks
go test ./...    # run tests
```

### Project layout

```
main.go               entrypoint → cmd.Execute()
cmd/root.go           Cobra root command and --copy flag
config/config.go      loads ~/standup.yaml + env
collector/            Activity type, Collector interface, git/GitHub collectors, concurrent runner
workday/              last-workday calculation (unit tested)
report/report.go      formats activities into the standup text
```
