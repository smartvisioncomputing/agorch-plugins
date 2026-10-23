# Slack

Connect your Agorch agents to Slack through Slack's official MCP server (`https://mcp.slack.com/mcp`). Agents can find a channel and post a message to it, such as a daily summary or an alert.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| Find a channel | Searches the workspace's public channels by name. | No |
| Send a message | Posts a message to a channel as the person who connected. | Yes |

## What connecting asks for

- **Your own Slack app.** Slack's MCP server only accepts apps that are internal to your workspace or listed in the Slack Marketplace, so your workspace creates one at `api.slack.com/apps`. Add the user scopes below and set the redirect URL to `https://app.agorch.ai/connection-oauth/callback`.
- **The app's client ID and client secret,** entered once in the Agorch connection. Agorch stores them; they never appear in a package or in this repository.
- **One sign-in.** Slack's MCP server works with user tokens, so messages are posted as the person who connects. The connection then belongs to your Agorch organization, and every agent in the organization uses it.

### Permissions

| Scope | In plain words |
|---|---|
| `search:read.public` | Find public channels |
| `chat:write` | Post messages as the person who connected |

## Examples

- "Post today's triage summary in #support-eng."
- "Tell #on-call that issue 412 was opened for the checkout error."

## Data sent to Slack

Channel search terms and message text.

## Good to know

- Sending a message changes data. Agorch never retries it on its own, and the starters ask a person to approve each call.
- Only the two tools above are allowed through.
