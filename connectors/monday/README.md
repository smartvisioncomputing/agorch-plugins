# monday.com

Connect your Agorch agents to monday.com through monday.com's official MCP server (`https://mcp.monday.com/mcp`). Agents can read boards, their columns and groups, and the items on them. This connector is read-only.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| Read a board | Reads a board's details, columns, groups and views. | No |
| List board items | Lists the items of a board, page by page, with optional filters. | No |

## What connecting asks for

- **One click and one sign-in.** Agorch already has a monday.com app, so there is nothing to create. Click Connect and authorize your monday.com account.
- The connection belongs to your Agorch organization, and every agent in the organization uses it.

### Permissions

| Scope | In plain words |
|---|---|
| `boards:read` | Read your boards, their columns and items |

## Examples

- "List the items in the Q4 Launch board that are still in Working on it."
- "Which columns does the Hiring board have?"

## Data sent to monday.com

Board IDs and item filters.

## Good to know

- Every tool call counts toward your monday.com account's daily API limit.
- Only the two read tools above are allowed through.
