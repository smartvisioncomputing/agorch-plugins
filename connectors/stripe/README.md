# Stripe

Connect your Agorch agents to Stripe through Stripe's official MCP server (`https://mcp.stripe.com/`). Agents can look up invoices and customers, create payment links and, with approval, issue refunds.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| Read Stripe data | Reads from the Stripe API: invoices, customers, payment links. | No |
| Write to Stripe | Creates, updates or deletes through the Stripe API: payment links, refunds. | Yes, and it can move money |

## What connecting asks for

- **A Stripe agent API key.** In the Stripe Dashboard, create a restricted key with the "Authorize agent access" option and give it only what your agents need: invoices, customers, payment links and refunds. From October 31, 2026, Stripe's MCP server refuses standard secret keys and restricted keys without agent access.
- **Paste the key once** into the Agorch connection. Agorch stores it; it never appears in a package or in this repository.
- The connection belongs to your Agorch organization, and every agent in the organization uses it.

### Permissions

There are no OAuth scopes: what the agents can do is exactly what you allow on the restricted key.

## Examples

- "List the invoices more than seven days overdue and the total still open."
- "Create a payment link for invoice in_1234 and send it to me for review."
- "Refund the duplicate charge on the Acme subscription." (asks for approval first)

## Data sent to Stripe

Stripe API requests, payment amounts and customer IDs.

## Good to know

- Writing to Stripe can move money, so every agent that uses it must ask a person to approve each call. Agorch never retries it on its own.
- Stripe adds its own approval step for refunds and payouts.
- Only the two tools above are allowed through.
