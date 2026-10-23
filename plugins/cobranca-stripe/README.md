# Stripe billing and recovery

**Money left on the table comes back.** The agent finds the overdue invoices, creates the payment link, records it in the CRM and leaves the reminder ready for your finance team to send.

## What it does

- Lists the overdue invoices in Stripe and adds up the open total.
- Uses the invoice's own link or creates a payment link. Only with your approval.
- Records a note and a follow-up task in HubSpot and prepares the reminder in Slack, following the overdue schedule.
- Refunds only on request and with two approvals: Agorch's and Stripe's own.

## Connect

- [Stripe](../../connectors/stripe/) (agent API key): read invoices and customers, create payment links and, with approval, refunds.
- [HubSpot](../../connectors/hubspot/): find the customer and record the collection.
- [Slack](../../connectors/slack/): send the reminder to the finance channel.

## Try it

- "List the invoices more than 7 days overdue, create a payment link for each one, record it in HubSpot and prepare the reminder in #finance."

## Approvals

Payment links, refunds, CRM writes and Slack messages always wait for a person to approve.

The agent's instructions, skill and test scenarios are written in Brazilian Portuguese.
