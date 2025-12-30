# Usage Guide

## Available Tools

The Warframe MCP server provides the following tools that can be used with Claude or any other MCP-compatible client:

### 1. `recommend_prime_farming`

Get intelligent recommendations for which prime parts to farm based on current world state, vault status, and your relic inventory.

**Example queries:**
- "Which prime parts should I farm now?"
- "What relics are worth running right now?"
- "What's currently unvaulted?"

**Returns:**
- Currently unvaulted primes from Varzia (Prime Resurgence)
- Best active fissures to run (sorted by tier)
- Your relic inventory (if Alecaframe is configured)
- Market value insights for prime sets

---

### 2. `recommend_standing_farming`

Get recommendations for which syndicate standing to farm based on active missions, world cycles, and current events.

**Example queries:**
- "What standing should I farm?"
- "Which syndicate missions are available?"
- "Is it day or night in Cetus?"

**Returns:**
- Active syndicate missions
- Current open-world cycles (Cetus, Orb Vallis, Cambion Drift, Zariman)
- Events that provide standing rewards
- Best missions for standing farming

---

### 3. `recommend_credit_farming`

Get the best methods to farm credits based on current events and active bonuses.

**Example queries:**
- "What's the best way to gain credits?"
- "How can I farm credits quickly?"
- "Are there any credit bonuses active?"

**Returns:**
- Best credit farming methods (Index, Profit-Taker, Railjack, etc.)
- Efficiency ratings for each method
- Active credit bonuses and events
- Requirements and expected rewards

---

### 4. `get_world_state`

Get the current Warframe world state including active alerts, invasions, fissures, sorties, and events.

**Example queries:**
- "What's happening in Warframe right now?"
- "Show me the current world state"
- "Are there any alerts or invasions?"

**Returns:**
- Active alerts, arbitrations, and sorties
- Void fissure count
- Invasion status
- Baro Ki'Teer status
- Open world cycles
- Current events

---

### 5. `get_baro_status`

Get information about Baro Ki'Teer (Void Trader) including his current location and inventory.

**Example queries:**
- "Is Baro here?"
- "What is Baro selling?"
- "When will Baro arrive?"

**Returns:**
- Current location (if visiting)
- Full inventory with ducats and credit costs
- Arrival/departure times
- Active status

---

### 6. `get_varzia_status`

Get information about Varzia and Prime Resurgence, including currently available primes and rotation schedule.

**Example queries:**
- "What's in the current Prime Resurgence rotation?"
- "When does the next rotation start?"
- "Which primes can I get from Varzia?"

**Returns:**
- Current rotation schedule
- Featured prime warframes and weapons
- Rotation start/end dates
- All items in the rotation

---

### 7. `get_market_price`

Get the current market price and trading information for a Warframe item from warframe.market.

**Parameters:**
- `item_name` (string): The name of the item to look up

**Example queries:**
- "What's the price of Rhino Prime Set?"
- "How much is Maiming Strike worth?"
- "Check the market price for Ember Prime Blueprint"

**Returns:**
- Average sell price
- Average buy price
- Minimum sell price / maximum buy price
- Number of active orders
- Trading recommendation

---

### 8. `search_market_items`

Search for tradeable items on warframe.market by name.

**Parameters:**
- `query` (string): Search query

**Example queries:**
- "Search for prime mods"
- "Find items related to Mesa"
- "Search for syndicate weapons"

**Returns:**
- List of matching items (up to 20 results)
- Item names for further price lookups

---

### 9. `get_my_relics`

Get your personal relic inventory from Alecaframe.

**Requirements:**
- Alecaframe credentials must be configured (`ALECAFRAME_USER_HASH` and `ALECAFRAME_PUBLIC_TOKEN`)

**Example queries:**
- "Show me my relic inventory"
- "What relics do I have?"
- "Do I have any vaulted relics?"

**Returns:**
- Complete list of your relics organized by era (Requiem, Axi, Neo, Meso, Lith)
- Total counts for each relic
- Refinement breakdown (Intact, Exceptional, Flawless, Radiant)
- Vault status for each relic

---

### 10. `get_my_stats`

Get your personal trading and account statistics from Alecaframe.

**Requirements:**
- Alecaframe credentials must be configured (`ALECAFRAME_USER_HASH` and `ALECAFRAME_PUBLIC_TOKEN`)

**Example queries:**
- "Show me my trading stats"
- "How much platinum have I earned?"
- "What are my account statistics?"

**Returns:**
- Total number of trades
- Total platinum earned
- Total platinum spent
- Net platinum (earned - spent)
- Other account metrics

---

## Example Conversations

### Finding what to farm

**You:** "I have some free time tonight. What should I farm in Warframe?"

**Claude (using the tools):**
*Uses `recommend_prime_farming`, `get_world_state`, and `recommend_credit_farming`*

"Based on the current world state, here are my recommendations:

**Prime Farming:**
- Mesa Prime, Ash Prime, Volt Prime, and Ember Prime are currently unvaulted through Prime Resurgence
- There are 8 active fissures, with good Axi missions available
- Best fissure: Axi Survival on Mot (Void) - ends in 45 minutes

**Credit Farming:**
- The Index is your best bet for pure credits (75k-200k per round)
- Alternatively, Profit-Taker gives 125k + credit booster if you have it unlocked

I'd recommend running some Axi fissures to get high-tier prime parts from the current rotation!"

---

### Checking Baro's inventory

**You:** "Is Baro here and what's he selling?"

**Claude (using the tools):**
*Uses `get_baro_status`*

"Yes! Baro Ki'Teer is currently at Tenno Relay on Earth. He brought 12 items this week, including:
- Primed Continuity (350 ducats, 110,000 credits)
- Prisma Grakata (400 ducats, 150,000 credits)
- [... more items]

He'll be here for another 1 day and 8 hours."

---

### Planning standing farming

**You:** "I need to farm some standing for Steel Meridian. What's the best approach?"

**Claude (using the tools):**
*Uses `recommend_standing_farming` and `get_world_state`*

"Here's your standing farming plan:

**Syndicate Missions:**
- Steel Meridian has 3 missions available today
- Locations: [mission nodes]

**Open World Bounties:**
- Plains of Eidolon is currently in Day cycle (2h 15m remaining)
- Great time to run Ostron bounties for standing

**Tip:** Complete your daily syndicate missions first for guaranteed standing, then supplement with bounties if needed!"

---

### Checking your account data

**You:** "Show me my relic inventory and trading stats"

**Claude (using the tools):**
*Uses `get_my_relics` and `get_my_stats`*

"Here's your account information:

**Relic Inventory:**
You have 47 unique relics across all eras:

- Axi: 15 relics (including 3 vaulted)
  - Axi A5: 8 total (2 Intact, 3 Exceptional, 2 Flawless, 1 Radiant) [VAULTED]
  - Axi M3: 12 total (5 Intact, 4 Exceptional, 2 Flawless, 1 Radiant)
  - [... more relics]

- Neo: 18 relics
- Meso: 10 relics
- Lith: 4 relics

**Trading Statistics:**
- Total Trades: 127
- Platinum Earned: 3,450
- Platinum Spent: 1,890
- Net Platinum: +1,560

You have a healthy relic collection! I'd recommend focusing on running those vaulted Axi A5 relics while they're still valuable."

---

## Configuration Tips

### Basic Setup (No Alecaframe)

You can use most features without Alecaframe configuration. The server will still provide:
- World state information
- Prime vault/resurgence data
- Market prices
- General recommendations

### Full Setup (With Alecaframe)

To get personalized recommendations based on your relic inventory:

1. In Warframe, open the Stats tab
2. Generate a public token with 'relic' access enabled
3. Copy your userHash (keep this private!)
4. Create `config.json`:

```json
{
  "alecaframe": {
    "userHash": "your-user-hash-here",
    "publicToken": "your-public-token-here"
  }
}
```

With this setup, the server can:
- Show your complete relic inventory with the `get_my_relics` tool
- Display your trading statistics with the `get_my_stats` tool
- Provide personalized recommendations in the `recommend_prime_farming` tool
- Help prioritize farming based on what you already have

---

## Tips for Best Results

1. **Ask specific questions**: The more specific your question, the better the recommendations
2. **Combine tools**: Claude can use multiple tools together to give you comprehensive advice
3. **Check regularly**: World state changes frequently, so recommendations update in real-time
4. **Market prices**: Remember that warframe.market prices fluctuate - check multiple times before trading

---

## Troubleshooting

### "No config file found" warning

This is normal if you haven't set up Alecaframe. All features except relic inventory will work.

### Market price not found

- Check the item name spelling
- Try searching with `search_market_items` first
- Some items may not be tradeable

### Rate limiting

The server automatically handles rate limiting for warframe.market (3 requests/second) and uses caching to minimize API calls.
