# Jira + Slack ticket triage

**No more tickets without an owner.** Every day the agent picks up the new Jira tickets, sets priority and component, suggests who should handle each one and posts a summary in your team's channel. It also flags the tickets that already missed the SLA.

## What it does

- Reads the new tickets with JQL (by default, the last 24 hours of the project).
- Sets the priority (Highest to Low) and the component, and suggests an assignee from the component's history.
- Writes the triage to the ticket and comments the reason. Only with your approval.
- Posts a summary in Slack and chases the tickets outside the SLA (4 hours for Highest, by default).

## Connect

- [Atlassian (Jira)](../../connectors/atlassian/): read, edit and comment on issues.
- [Slack](../../connectors/slack/): find the channel and post the summary.

## Try it

- "Triage the tickets created in the last 24 hours in project SUP and post the summary in #support-eng."
- "Which Highest tickets in SUP have been open for more than 4 hours without an assignee?"

## Approvals

Editing a ticket, commenting and posting in Slack always wait for a person to approve.

The agent's instructions, skill and test scenarios are written in Brazilian Portuguese.
