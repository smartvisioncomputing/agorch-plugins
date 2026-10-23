# Sentry to GitHub on-call

**A new error becomes an issue with an owner in minutes.** The agent groups the new Sentry errors, opens the GitHub issue with the stack trace and release, tags the code owner and notifies the on-call channel.

## What it does

- Finds the new errors and groups them by exception and file.
- Never duplicates: it looks for an existing issue before opening one.
- Finds the owner from `CODEOWNERS` or from the file's commits.
- Opens the issue and notifies Slack. Only with your approval.

## Connect

- [Sentry](../../connectors/sentry/): read the new errors with stack trace, release and event count.
- [GitHub](../../connectors/github/) (your organization's GitHub App): avoid duplicates, find the owner and open the issue.
- [Slack](../../connectors/slack/): notify the on-call channel with the issue link.

## Try it

- "Check the new errors from the last hour in Sentry's web project, open an issue in acme/web for the ones that have none, and notify #on-call."

## Approvals

Opening an issue and posting in Slack always wait for a person to approve.

The agent's instructions, skill and test scenarios are written in Brazilian Portuguese.
