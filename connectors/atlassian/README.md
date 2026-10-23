# Atlassian (Jira)

Connect your Agorch agents to Jira through Atlassian's official Rovo MCP server (`https://mcp.atlassian.com/v2/mcp`). Agents can find your Atlassian sites, search issues with JQL, read, create and edit issues, and comment on them.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| List Atlassian sites | Lists the sites the connected account can open. | No |
| Search issues (JQL) | Finds issues with a JQL query. | No |
| Read an issue | Reads one issue and its fields. | No |
| Create an issue | Opens an issue in a project: summary, type, description and project fields. | Yes |
| Edit an issue | Changes editable fields such as priority, component and assignee. | Yes |
| Comment on an issue | Adds or edits a comment. | Yes |

## What connecting asks for

- **One click and one sign-in.** Agorch already has an Atlassian app, so there is nothing to create. Click Connect, sign in to Atlassian and pick the site.
- **Rovo MCP turned on for your site.** A site admin enables it in Atlassian administration. Calls made through Rovo MCP count toward your Rovo credits.
- The connection belongs to your Agorch organization, and every agent in the organization uses it.

### Permissions

| Scope | In plain words |
|---|---|
| `search:jira:agent-interface` | Search Jira issues with JQL |
| `read:jira:agent-interface` | Read Jira issues and the sites you can access |
| `write:jira:agent-interface` | Create, edit and comment on Jira issues |

## Examples

- "Find the tickets created in project SUP in the last 24 hours."
- "Set SUP-142 to High priority, component Billing, and explain why in a comment."
- "Open a bug in SUP for the checkout timeout with the steps from this thread."

## Data sent to Atlassian

JQL queries, issue fields and comment text.

## Good to know

- Creating, editing and commenting change data. Agorch never retries them on its own, and the starters ask a person to approve each call.
- Only the six tools above are allowed through. Tools that Rovo MCP hides by default stay blocked.
