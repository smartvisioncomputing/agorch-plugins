# SDR with RD Station

**No lead goes cold, on the CRM Brazilian sales teams use.** The agent qualifies every new RD Station conversion, moves the good ones through the funnel, creates the follow-up and lets the salesperson know in Slack.

## What it does

- Finds the new conversions and scores them from 0 to 3 (fit, contact data and interest).
- Updates the contact, moves the deal and creates the follow-up task. Only with your approval.
- Tells each salesperson about their leads in Slack, without exposing email addresses or phone numbers.

## Connect

- [RD Station](../../connectors/rdstation/) (CRM or Marketing, beta): read leads, record the qualification, move the deal and create the task.
- [Slack](../../connectors/slack/): notify the salesperson.

## Try it

- "Qualify today's conversions in RD Station, move the qualified ones in the funnel, create the follow-up and tell the salesperson in #sales."

## Approvals

Every CRM write and every Slack message waits for a person to approve.

Using HubSpot? See [SDR with HubSpot](../sdr-hubspot/).

The agent's instructions, skill and test scenarios are written in Brazilian Portuguese.
