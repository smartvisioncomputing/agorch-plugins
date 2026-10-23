# GitHub

Connect your Agorch agents to GitHub through GitHub's official remote MCP server (`https://api.githubcopilot.com/mcp/`). Agents can search issues, read files such as `CODEOWNERS`, list commits, and open or update issues.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| Search issues | Finds issues with GitHub's search syntax. | No |
| Read a file | Reads a file from a repository, for example `CODEOWNERS`. | No |
| List commits | Lists the recent commits of a repository. | No |
| Open or update an issue | Creates an issue or updates an existing one. It never deletes. | Yes |

## What connecting asks for

- **Your own GitHub App.** GitHub does not let other apps register themselves, so your organization creates a GitHub App and keeps control of it. Give it these permissions: Issues (read and write), Contents (read) and Metadata (read). Set the callback URL to `https://app.agorch.ai/connection-oauth/callback`, then install the app on the repositories your agents may use.
- **The app's client ID and client secret,** entered once in the Agorch connection. Agorch stores them; they never appear in a package or in this repository.
- **One sign-in.** The person who connects authorizes the app. The connection then belongs to your Agorch organization, and every agent in the organization uses it.

### Permissions

| Scope | In plain words |
|---|---|
| `repo` | Read and write issues and read code in your repositories |

A GitHub App ignores OAuth scopes: what the app may do is set by its own permissions above. The `repo` scope only matters if you connect with a classic OAuth App instead.

## Examples

- "Look for an open issue about the checkout timeout before you open a new one."
- "Read `CODEOWNERS` in `acme/web` and tell me who owns `src/payments/`."
- "List the last ten commits in `acme/web`."

## Data sent to GitHub

Repository names, issue titles and bodies, and file paths.

## Good to know

- Opening or updating an issue changes data. Agorch never retries it on its own, and the starters ask a person to approve each call.
- Only the four tools above are allowed through. Everything else the GitHub MCP server offers stays blocked.
