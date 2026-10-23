# SDR with HubSpot

**No lead goes cold.** The agent qualifies every new HubSpot contact, hands the good ones to the right salesperson with a follow-up task already created, and lets them know in Slack.

## What it does

- Finds the new contacts and scores them from 0 to 3 (fit, contact data and interest), using only CRM data.
- Moves the qualified ones through the funnel and creates the follow-up task for the owner. Only with your approval.
- Tells each salesperson about their leads in Slack, without exposing email addresses or phone numbers.

## Connect

- [HubSpot](../../connectors/hubspot/): read leads and deals, record the qualification and the task.
- [Slack](../../connectors/slack/): notify the salesperson.

## Try it

- "Qualify the contacts created today in HubSpot, create the follow-up for the qualified ones and tell the salesperson in #sales."

## Approvals

Every CRM write and every Slack message waits for a person to approve.

Using RD Station? See [SDR with RD Station](../sdr-rdstation/).

The agent's instructions, skill and test scenarios are written in Brazilian Portuguese.
