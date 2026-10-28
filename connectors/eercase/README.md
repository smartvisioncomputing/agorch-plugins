# EERCase

Connect your Agorch agents to [EERCase](https://eercase.cin.ufpe.br), the EER modeling tool of CIn/UFPE, through its MCP server (`https://eercase.cin.ufpe.br/mcp`). Agents read, validate and render the diagrams the signed-in user can see, generate SQL, and propose or apply changes.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| List projects | Projects you own or that were shared with you, with your permission and version. | No |
| Read a diagram | The diagram as compact text, with every element id and the project version. | No |
| Render a diagram | A PNG of the saved diagram. | No |
| Validate a diagram | The editor's own validation; can simulate changes first. | No |
| Generate SQL | The relational SQL, as the editor exports it. | No |
| Modeling guide | The EER metamodel and the read → validate → propose workflow. | No |
| Propose changes | Shows an Apply / Discard card in the editor the user has open and waits up to 45 s for the answer. | Only if the user applies |
| Read a proposal | The answer to a proposal that was still pending. | No |
| Apply changes | Saves at once; refused with `version_conflict` if the project changed since it was read. | Yes |
| Create a project | A new empty project. | Yes |

## What connecting asks for

- **One sign-in with Google at EERCase, no app to create.** EERCase has no client registration: Agorch identifies itself with its public client document (`https://app.agorch.ai/oauth/client-metadata.json`), and EERCase shows "Agorch" on its consent page.
- The connection belongs to your Agorch organization, and every agent in the organization uses it, as the EERCase user who connected.
- If that user is not cleared for AI tools in EERCase, the sign-in stops at EERCase's "not authorized" page, and a tool answers `ai_not_authorized`: "this user is not authorized to use AI tools in EERCase". Ask the EERCase team to authorize the user, then connect again.

### Permissions

| Scope | In plain words |
|---|---|
| `eercase:read` | List, read, render and validate your diagrams and generate SQL |
| `eercase:write` | Propose and apply changes to your diagrams (you can uncheck it at the consent) |

## Examples

- "List my EERCase projects and validate the library one."
- "Model a minimundo of a clinic and propose it in the editor I have open."

## Data sent to EERCase

Project ids, diagram operations and the text of new element names. What an agent reads (names, structure, image, SQL) goes to the agent's model provider.

## Good to know

- Limits per user: 60 reads and 10 writes per minute.
- Revoke at EERCase (Settings → Connected apps) or disconnect in Agorch.
- Not in the marketplace index: a sealed package (ComponentPackage v1.1) cannot carry a public client without registration yet, so this stack is applied to an organization with `aicl sync`.
