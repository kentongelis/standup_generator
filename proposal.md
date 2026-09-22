# Project Proposal: `standup`

**Author:** Kenton
**Language:** Go
**Target platform:** GitHub

## Summary

`standup` is a command-line tool that generates a daily standup update automatically. It gathers a developer's recent activity from local git repositories and GitHub (commits, pull requests, reviews, and issues), runs those lookups concurrently, and produces a clean **Yesterday / Today / Blockers** summary that can be copied to the clipboard or posted directly to Slack.

## Problem

Every morning, each engineer spends time reconstructing what they did the previous day: scrolling through commit history, checking which PRs they reviewed, and remembering which issues they touched. This is repetitive, easy to get wrong, and tends to produce vague updates ("worked on the API stuff"). The information already exists in git and GitHub; nobody should have to assemble it by hand.

## Proposed Solution

A single command:

```bash
standup
```

produces output like:

```
Yesterday
  • Merged PR #142: Add rate limiting to /auth endpoints
  • Reviewed PR #139: Refactor user service (approved)
  • Closed issue #88: Login redirect loop on Safari

Today
  • Continue work on branch feat/session-refresh (3 commits, no PR yet)
  • Address review comments on PR #145

Blockers
  • PR #145 has been waiting on review for 2 days
  • CI failing on PR #147
```

## Features

### MVP (must ship this sprint)

- **Git collector:** reads configured local repositories for the user's commits since the last workday.
- **GitHub collector:** fetches PRs opened, merged, reviewed, and commented on, plus issues opened or closed, using the GitHub API.
- **Concurrent collection:** all collectors run in parallel so the tool stays fast across many repos.
- **"Last workday" logic:** on Monday, the report covers Friday rather than Sunday.
- **Template output:** a deterministic, offline-friendly summary in Yesterday / Today / Blockers format.
- **Clipboard copy:** `--copy` places the summary on the clipboard for pasting into Slack.
- **Config file:** `~/.standup.yaml` stores the GitHub token, username, and repo list.

### Next (if time allows)

- **AI summaries:** an `--ai` flag sends the normalized activity to an LLM to produce a more natural, concise write-up. The template output remains the fallback if no API key is set.
- **Slack posting:** `--post` sends the summary to a channel through an incoming webhook.
- **Blocker detection:** flags PRs waiting on review longer than a threshold and PRs with failing checks.

### Stretch

- **Team digest:** `standup team` aggregates summaries for a list of teammates into one message.
- **Sprint recap:** `standup --since 2w` generates a sprint-length summary suitable for retrospectives.

## Architecture

```
            ┌────────────────┐
            │   Cobra CLI    │
            └───────┬────────┘
                    │ since = lastWorkday()
        ┌───────────┴───────────┐
        ▼                       ▼
┌───────────────┐       ┌────────────────┐
│ Git collector │       │GitHub collector│   (goroutines via errgroup)
└───────┬───────┘       └───────┬────────┘
        └───────────┬───────────┘
                    ▼
            ┌───────────────┐
            │  Normalizer   │  merge, dedupe, group commits under PRs
            └───────┬───────┘
                    ▼
            ┌───────────────┐
            │  Summarizer   │  template or LLM
            └───────┬───────┘
                    ▼
        terminal / clipboard / Slack
```

Each data source implements a shared interface, which keeps collectors independent and makes new sources easy to add later:

```go
type Activity struct {
    Source    string    // "git", "github"
    Kind      string    // "commit", "pr_merged", "pr_review", "issue_closed", ...
    Title     string
    URL       string
    Repo      string
    Timestamp time.Time
}

type Collector interface {
    Name() string
    Collect(ctx context.Context, since time.Time) ([]Activity, error)
}
```

If one collector fails (for example, a network error), the others still report, and the output notes which source was unavailable.

## Tech Stack

| Purpose | Package |
|---|---|
| CLI framework | `github.com/spf13/cobra` |
| Configuration | `github.com/spf13/viper` |
| GitHub API | `github.com/google/go-github` |
| Local git history | `github.com/go-git/go-git` |
| Concurrency | `golang.org/x/sync/errgroup` |
| Clipboard | `github.com/atotto/clipboard` |
| Terminal spinner/styling | `github.com/charmbracelet/lipgloss` (optional) |

## Why Go

- **Goroutines** make it natural to query many repositories and API endpoints at once, keeping runtime to a few seconds.
- **Single static binary** means teammates can install it with `go install` or a downloaded release, with no runtime dependencies.
- **Mature ecosystem** for CLIs (Cobra, Viper) and the GitHub API (go-github) means effort goes into the product rather than plumbing.

## Sprint Plan

| Phase | Work |
|---|---|
| Days 1–2 | Project scaffold, Cobra commands, config loading, `Activity` and `Collector` types |
| Days 3–4 | Git collector and GitHub collector, last-workday logic |
| Days 5–6 | Concurrent orchestration, normalizer, template output, clipboard |
| Days 7–8 | Blocker detection, Slack posting, `--ai` summaries |
| Days 9–10 | Tests, README, release binaries, demo preparation |

## Success Criteria

- Generates an accurate standup in under 5 seconds for a user with 10+ repositories.
- Works with only a GitHub token configured (no LLM or Slack required).
- A new user can install and produce their first standup in under 5 minutes by following the README.
- Positive feedback and adoption from teammates via the #general feedback form.

## Risks and Mitigations

| Risk | Mitigation |
|---|---|
| GitHub API rate limits | Use authenticated requests, the Search API for batched queries, and cache results for the day |
| Token security | Read the token from an environment variable or config file with restricted permissions; never log it |
| LLM output inaccuracies | AI mode is opt-in, and the summary is generated only from collected activity, with links to each item |
| Scope creep | MVP features are fixed; everything else is explicitly labeled Next or Stretch |

## Demo Plan

1. Open the retrospective presentation with a standup the tool generated about building itself.
2. Run `standup` live, with a spinner showing each collector finishing concurrently.
3. Show the `--ai` version side by side with the template version.
4. Post the result to Slack with `--post` to close the loop.
