# Notion

Connect your Agorch agents to Notion through Notion's official MCP server (`https://mcp.notion.com/mcp`). Agents can search your workspace and read pages and databases. This connector is read-only.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| Search Notion | Searches the pages and databases you shared. | No |
| Read a page or database | Reads the content of a page or the rows of a database. | No |

## What connecting asks for

- **One click and one sign-in.** Agorch already has a Notion app, so there is nothing to create. Click Connect, sign in to Notion and choose the pages the agents may see.
- The connection belongs to your Agorch organization, and every agent in the organization uses it.

### Permissions

| Scope | In plain words |
|---|---|
| `default` | Search and read the pages you share during the connection |

Notion has no finer scopes. What the agents can see is exactly the set of pages you pick on Notion's consent screen, and you can change it later in Notion.

## Examples

- "Find our refund policy page and summarize it."
- "Read the Roadmap database and list the items due this month."

## Data sent to Notion

Search queries and page IDs.

## Good to know

Only the two read tools above are allowed through.
