# Nuvemshop store operator

**Your store checked every day, without opening the admin panel.** The agent sends the sales summary, warns you when stock is running low and handles prices and coupons with your approval.

## What it does

- Daily summary in Slack: paid orders, revenue, average ticket and best sellers.
- Low-stock alert (5 units or fewer, by default).
- Changes prices and creates coupons only with your approval, and refuses a discount above 30% without confirmation.

## Connect

- [Nuvemshop](../../connectors/nuvemshop/) (beta): read orders and stock, change prices and create coupons.
- [Slack](../../connectors/slack/): post the summary and the alerts in the store's channel.

## Try it

- "Summarize yesterday's orders (count, revenue and average ticket) and post it in #store."

## Approvals

Price changes, coupons and Slack messages always wait for a person to approve.

The agent's instructions, skill and test scenarios are written in Brazilian Portuguese.
