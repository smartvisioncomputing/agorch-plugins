# HubSpot

Connect your Agorch agents to your HubSpot CRM through HubSpot's official remote MCP server (`https://mcp.hubspot.com/`). Agents can search and read contacts, companies and deals, record notes and tasks, and find the salesperson who owns a record.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| Search CRM records | Finds contacts, companies and deals by filters. | No |
| Read CRM records | Reads the fields of one or more records. | No |
| Write CRM records | Creates or updates contacts, deals, notes and tasks. | Yes |
| Search owners | Finds the HubSpot users who own records, such as the salesperson of a lead. | No |

## What connecting asks for

- **Your own HubSpot MCP auth app.** HubSpot asks each customer to create one: go to `app.hubspot.com/l/mcp-auth-apps/`, create the app and set the redirect URL to `https://app.agorch.ai/connection-oauth/callback`.
- **The app's client ID and client secret,** entered once in the Agorch connection. Agorch stores them; they never appear in a package or in this repository.
- **One sign-in.** The person who connects authorizes HubSpot for your Agorch organization, and every agent in the organization uses that connection.

### Permissions

The connector asks for no OAuth scopes. HubSpot decides what the MCP server can reach from the tools installed on the auth app, so review that list when you create it.

## Examples

- "Find the contacts created today and show their company and lifecycle stage."
- "Add a note to the Acme deal saying the proposal was sent, and create a follow-up task for Friday."
- "Who owns the contact maria@example.com?"

## Data sent to HubSpot

CRM search filters, record fields and note text.

## Good to know

- Writing CRM records changes data. Agorch never retries it on its own, and the starters ask a person to approve each call.
- Only the four tools above are allowed through.
