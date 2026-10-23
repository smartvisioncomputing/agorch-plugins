# RD Station

Connect your Agorch agents to RD Station CRM and RD Station Marketing through RD Station's official MCP server (`https://mcp.rdstationmentor.com/`). Agents can find contacts and leads, update them, move deals to another stage and create follow-up tasks.

**Beta.** RD Station does not yet publish the server's OAuth details or tool names. The endpoint, scopes and tool names are confirmed with the first live connections during the beta. If a name differs, only the mapping changes: the tools your agents use keep their names in Agorch.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| Search contacts and leads | Finds contacts and recent conversions. | No |
| Update a contact | Changes a contact's fields, such as the qualification. | Yes |
| Move a deal | Moves a deal to another stage of the funnel. | Yes |
| Create a follow-up task | Creates a task for the deal's owner. | Yes |

## What connecting asks for

- **An RD Station CRM or RD Station Marketing subscription.**
- **One click and one sign-in.** Agorch registers the app with RD Station, so there is nothing to create. Click Connect and authorize your RD Station account.
- The connection belongs to your Agorch organization, and every agent in the organization uses it.

### Permissions

| Scope | In plain words |
|---|---|
| `crm` | Read and update contacts, deals and tasks in RD Station CRM |
| `marketing` | Read leads and conversions in RD Station Marketing |

## Examples

- "Show today's conversions with their source and company."
- "Move the Acme deal to Proposal and create a follow-up task for tomorrow."

## Data sent to RD Station

Contact fields, deal stages and task text.

## Good to know

- Updating contacts, moving deals and creating tasks change data. Agorch never retries them on its own, and the starters ask a person to approve each call.
- Only the four tools above are allowed through.
