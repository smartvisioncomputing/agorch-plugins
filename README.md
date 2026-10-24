# Agorch plugins

The first-party packages of the [Agorch](https://agorch.ai) marketplace: 10 connectors to official vendor MCP servers and 6 ready-made starter agents. Everything is written in AICL 1.1 (`plugin.aicl` plus a `stack/`) and sealed into content-addressed packages by `aicl plugin build`.

## Starter agents

| Starter | What it does | Connects |
|---|---|---|
| [Jira + Slack ticket triage](plugins/triagem-jira-slack/) | Sorts new tickets by priority and component, suggests an owner, posts a summary and chases the SLA. | Atlassian, Slack |
| [SDR with HubSpot](plugins/sdr-hubspot/) | Qualifies new leads, moves them through the funnel, creates the follow-up and alerts the salesperson. | HubSpot, Slack |
| [SDR with RD Station](plugins/sdr-rdstation/) | The same SDR, on the CRM Brazilian sales teams use. | RD Station, Slack |
| [Stripe billing and recovery](plugins/cobranca-stripe/) | Overdue invoices, payment links, CRM notes and reminders; refunds only with approval. | Stripe, HubSpot, Slack |
| [Nuvemshop store operator](plugins/operador-nuvemshop/) | Daily sales summary, low-stock alerts, price and coupon changes with approval. | Nuvemshop, Slack |
| [Sentry to GitHub on-call](plugins/plantao-sentry-github/) | A new error becomes an issue with the stack trace, release and code owner, and the channel is notified. | Sentry, GitHub, Slack |

Every tool call that changes something (editing, commenting, posting, writing to the CRM, payment links, refunds, prices, coupons, opening issues) waits for a person to approve it (`approval = "human_required"`). The starters are built for Brazilian teams: their instructions and scenarios are in Brazilian Portuguese.

## Connectors

Each connector reaches the vendor's own MCP server and allows only the tools listed in its README.

| Connector | What agents can do | Category | Connection | Tools |
|---|---|---|---|---|
| [HubSpot](connectors/hubspot/) | Search, read and update contacts and deals. | Sales & Marketing | OAuth, your own HubSpot MCP auth app | 4 |
| [RD Station](connectors/rdstation/) (beta) | Leads, funnel and tasks. | Sales & Marketing | OAuth | 4 |
| [Stripe](connectors/stripe/) | Invoices, customers and payment links; refunds with approval. | Payments & Commerce | Agent API key | 2 |
| [Nuvemshop](connectors/nuvemshop/) (beta) | Orders, stock, prices and coupons. | Payments & Commerce | OAuth | 4 |
| [Atlassian](connectors/atlassian/) | Search, create, edit and comment on Jira issues. | Work & Docs | OAuth | 6 |
| [monday.com](connectors/monday/) | Read boards and items. | Work & Docs | OAuth | 2 |
| [Notion](connectors/notion/) | Search and read pages and databases. | Work & Docs | OAuth | 2 |
| [Slack](connectors/slack/) | Find a channel and post a message. | Communication | OAuth, your own Slack app | 2 |
| [GitHub](connectors/github/) | Search and open issues, read `CODEOWNERS` and commits. | Engineering | OAuth, your own GitHub App | 4 |
| [Sentry](connectors/sentry/) | Errors, stack traces and releases, read-only. | Engineering | OAuth | 2 |

Every connection belongs to an Agorch organization, not to a single user. Credentials are entered in Agorch at connection time and never live in this repository.

## Layout

```
connectors/<vendor>/   plugin.aicl (version, license, readme, release notes, scope labels, source), stack/, README.md
plugins/<starter>/     plugin.aicl (requires, permissions), stack/ (agent, skill in SKILL.md, scenarios), aicl.lock.json, README.md
index/                 one sealed package per connector (ao.<vendor>-mcp@<version>.json)
starters/              each starter's sealed agent_system package (ao.<starter>@<version>.json)
catalog.json           marketplace listing: name, summary, category, rank, badges and auth kind per connector, plus the starters
cmd/publish/           signs index/ and publishes it to an Agorch registry
scripts/check.sh       everything CI checks
scripts/index.sh       regenerates index/, the locks and starters/
scripts/publish.sh     publishes index/ and catalog.json to a registry (idempotent)
```

Each connector becomes the package `ao/<vendor>-mcp`. A starter asks for it with `requires "<vendor>" { package = "ao/<vendor>-mcp", version = "^0.1.0" }`, and its committed `aicl.lock.json` pins the version, the digest and the `side_effect` of every tool it uses.

Each skill is a `SKILL.md` in the [Agent Skills](https://agentskills.io/specification) format: frontmatter with `name` (the skill id in `stack.aicl` and the folder name) and `description`, then the body.

## What a sealed connector carries

Besides the tools and the connection definition, every package in `index/` seals its presentation (`package.presentation`): the license, the README above, the release notes, a plain-language label for each OAuth scope, and the source it was built from (`repository`, `commit`, `tree` and `path`). The commit is the last one that touched `connectors/<vendor>/`, and the tree is that folder's git tree at that commit, so anyone can check that a package was built from exactly the code in this repository.

## Checking

```sh
# Static checks, no compiler needed (only git and jq): catalog.json x index/ x plugins/ and the secret scan.
scripts/check.sh static

# Full check with the aicl compiler: fmt -check, plugin validate, every connector in index/ rebuilt
# from the commit it names and compared byte for byte, starters built twice and compared with starters/.
AICL=/path/to/aicl scripts/check.sh
```

Every pull request runs both. The `aicl` compiler is not published yet, so the full check needs the `AICL_REPO_TOKEN` secret, which a pull request from a fork does not get: such a pull request fails when it changes `index/`, `starters/`, `connectors/`, `plugins/`, `scripts/`, `.github/` or `catalog.json`, and a maintainer pushes it to a branch of this repository to run the check.

## Releasing

1. Change the connector or the starter. For a connector, raise `version` in both `plugin.aicl` and `stack/stack.aicl` and update `release_notes`: a published version is never sealed twice.
2. Commit the change.
3. Run `AICL=/path/to/aicl scripts/index.sh` on the clean tree. Phase 1 rebuilds `index/`, pinning each connector to its own commit; phase 2 resolves the locks against the new `index/` and rebuilds `starters/`.
4. Commit `index/`, the locks and `starters/` together, and run `scripts/check.sh`.

Release tags follow `connectors/<id>@<version>` or `plugins/<id>@<version>`.

## Publishing

`scripts/publish.sh` takes `index/` and then `starters/` to an Agorch registry through the ordinary publisher path: publisher account `publisher:ao`, the `first_party` claim on the `ao` namespace, a signing key, the signed artifact, listing, release, automated checks, platform review and the `stable` channel. It then applies the listing fields from `catalog.json` (`category`, `rank`, `vendor_official`, `beta`). Every step reads the registry's state first, so a second run changes nothing. A starter is published only after every connector it requires is on the stable channel; otherwise it is skipped and reported (the run still publishes the others). Starters get no category, rank or badges from `catalog.json`: the API defaults apply.

```sh
export AGORCH_API=https://<api>                         # API origin, no trailing slash
export AGORCH_PUBLISHER_TOKEN=<id token>                # a member of publisher:ao
export AGORCH_PLATFORM_TOKEN=<id token>                 # a platform reviewer, not the publisher
export AGORCH_PUBLISHER_SEED_FILE=<file with 64 hex>    # Ed25519 seed of the signing key
export AGORCH_OWNER_EMAIL=<email>                       # first run only: the account owner
AICL=/path/to/aicl scripts/publish.sh
```

`scripts/publish.sh` runs the full check first and publishes nothing if any package in `index/` is not what its source builds.

Environment values are never printed. The publisher and the reviewer must be different people: the registry refuses a self-approval.

## Security

- No literal secrets. Credentials enter only through `secret_ref` and their values stay in Agorch. `aicl` scans every manifest, README and package, and `scripts/check.sh` scans the repository without printing what it finds.
- Minimal scopes per connector, and an allowlist with only the tools the starters use.
- Writes are never retried automatically, and every starter asks a person before any change.
- Report a vulnerability privately to contact@agorch.ai.

## License

[Apache-2.0](LICENSE). Vendor names belong to their owners; a connector labeled official uses the vendor's own MCP server.
