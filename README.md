# standup_generator

A CLI tool that generates a daily standup update (Yesterday / Today / Blockers) from your git and GitHub activity. See `proposal.md` for the full design.

## Status

Early development. Working so far:

- Cobra CLI entrypoint (`standup`)
- Config loading from `~/.standup.yaml` (with `GITHUB_TOKEN` env override)
- `Activity` / `Collector` interface, with a `Fake` collector wired in for end-to-end testing

Not yet implemented: real git/GitHub collectors, last-workday logic (currently hardcoded to the last 24 hours), normalizer, summarizer, clipboard/Slack output.

## Usage

```bash
go run .
```

Currently prints activity from the fake collector; real collectors are still in progress.
