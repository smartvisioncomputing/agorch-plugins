# Sentry

Connect your Agorch agents to Sentry through Sentry's official MCP server (`https://mcp.sentry.dev/mcp`). Agents can search error issues and read an issue, event or trace with its stack trace, release and tags. This connector is read-only.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| Search issues | Finds error issues with Sentry's search syntax, for example the new ones in the last hour. | No |
| Read an issue or event | Reads an issue, event or trace: stack trace, release, tags and event count. | No |

## What connecting asks for

- **One click and one sign-in.** Agorch already has a Sentry app, so there is nothing to create. Click Connect and authorize your Sentry organization.
- The connection belongs to your Agorch organization, and every agent in the organization uses it.

### Permissions

| Scope | In plain words |
|---|---|
| `org:read` | Read your Sentry organization |
| `project:read` | Read your projects and their releases |
| `event:read` | Read issues, events and stack traces |

## Examples

- "What new errors did the web project get in the last hour?"
- "Show the stack trace and release of issue WEB-3F2."

## Data sent to Sentry

Search queries and issue IDs.

## Good to know

Only the two read tools above are allowed through.
