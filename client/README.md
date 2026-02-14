# Reward Chain Client

A clean and slim Node.js client for interacting with the Reward Chain Cosmos SDK blockchain.

## Features

- ✅ Create partners on-chain
- ✅ List all partners
- ✅ Get partner by ID
- ✅ Clean and minimal API
- ✅ Uses generated proto files

## Installation

```bash
cd client
npm install
```

## Configuration

Set environment variables or modify `example.js`:

```bash
export RPC_ENDPOINT="http://localhost:26657"
export MNEMONIC="your mnemonic phrase here"
export ADDRESS_PREFIX="reward"
export GAS_PRICE="0.0001stake"  # Optional, defaults to 0.0001stake
```

## Usage

### Basic Example

```javascript
import RewardChainClient from "./client.js";

// Connect to the chain
const client = await RewardChainClient.connect(
  "http://localhost:26657",
  "your mnemonic phrase",
  "reward",           // address prefix
  "0.0001stake"      // gas price (optional)
);

// Create a partner
const result = await client.createPartner({
  name: "My Partner",
  category: "Retail",
  country: "US",
  currency: "USD",
  earnCostPerPoint: "1.0",
  burnCostPerPoint: "0.9",
  totalLiquidity: "10000",
});

console.log("Partner ID:", result.partnerId);
console.log("Transaction Hash:", result.transactionHash);

// List all partners
const partners = await client.listPartners();
console.log("Partners:", partners);

// Get a specific partner
const partner = await client.getPartner(1);
console.log("Partner:", partner);

// Disconnect
await client.disconnect();
```

### Run Example

```bash
npm start
```

Or:

```bash
node example.js
```

## API Reference

### `RewardChainClient.connect(rpcEndpoint, mnemonic, prefix, gasPrice)`

Creates a new client instance with signing capabilities.

**Parameters:**
- `rpcEndpoint` (string): RPC endpoint URL (e.g., "http://localhost:26657")
- `mnemonic` (string): Mnemonic phrase for the wallet
- `prefix` (string, optional): Bech32 address prefix (default: "reward")
- `gasPrice` (string, optional): Gas price (default: "0.0001stake")

**Returns:** `Promise<RewardChainClient>`

### `client.createPartner(partnerData, options)`

Creates a new partner on-chain.

**Parameters:**
- `partnerData` (object):
  - `name` (string): Partner name
  - `category` (string): Partner category
  - `country` (string): Partner country
  - `currency` (string): Partner currency
  - `earnCostPerPoint` (string): Earn cost per point
  - `burnCostPerPoint` (string): Burn cost per point
  - `totalLiquidity` (string): Total liquidity
- `options` (object, optional):
  - `memo` (string): Transaction memo
  - `fee` (Fee | string): Transaction fee (default: calculated from gas price)
  - `gas` (number): Gas limit (default: 200000)

**Returns:** `Promise<Object>` with `transactionHash`, `partnerId`, `height`, and `gasUsed`

### `client.listPartners(options)`

Lists all partners.

**Parameters:**
- `options` (object, optional):
  - `includeDisabled` (boolean): Include disabled partners (default: false)
  - `pagination` (object): Pagination options

**Returns:** `Promise<Array>` of partner objects

### `client.getPartner(partnerId)`

Gets a specific partner by ID.

**Parameters:**
- `partnerId` (number): Partner ID

**Returns:** `Promise<Object>` partner data

## Notes

- The client uses the generated proto files from the `proto/` directory
- All string amounts should be passed as strings (e.g., "1000" not 1000)
- The client automatically handles transaction signing and broadcasting
- Query operations use REST endpoints (port 1317) while transactions use RPC (port 26657)
- Gas price is configured when connecting the client (default: "0.0001stake")
- If fee is not provided, it's automatically calculated from the gas price and gas limit