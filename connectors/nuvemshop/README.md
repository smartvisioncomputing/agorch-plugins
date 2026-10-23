# Nuvemshop

Connect your Agorch agents to your Nuvemshop store through Nuvemshop's official admin MCP server (`https://admin-mcp.tiendanube.com/`, the host Nuvemshop's own metadata announces). Agents can read orders, products and stock, change prices and stock, and create coupons.

**Beta.** Nuvemshop does not publish the names of its MCP tools. They are confirmed against the live server during the beta. If a name differs, only the mapping changes: the tools your agents use keep their names in Agorch.

## What your agents can do

| Tool | What it does | Changes data |
|---|---|---|
| List orders | Lists the store's orders, with filters. | No |
| List products and stock | Lists products, prices and stock levels. | No |
| Update a product | Changes a product's price or stock. | Yes |
| Create a coupon | Creates a discount coupon with its rules. | Yes |

## What connecting asks for

- **One click and one sign-in.** Agorch already has a Nuvemshop app, so there is nothing to create. Click Connect, sign in as the store owner and authorize the store.
- The connection belongs to your Agorch organization, and every agent in the organization uses it.

### Permissions

| Scope | In plain words |
|---|---|
| `read_orders` | Read your store's orders |
| `read_products` | Read your products and stock |
| `write_products` | Change product prices and stock |
| `write_coupons` | Create discount coupons |

## Examples

- "Summarize yesterday's paid orders: count, revenue and average ticket."
- "Which products have five units or fewer in stock?"
- "Create a 10% coupon for the winter collection, valid until Sunday."

## Data sent to Nuvemshop

Order filters, product prices and stock, and coupon rules.

## Good to know

- Updating products and creating coupons change data. Agorch never retries them on its own, and the starters ask a person to approve each call.
- Only the four tools above are allowed through.
