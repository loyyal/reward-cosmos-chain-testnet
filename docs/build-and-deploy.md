# Build & Deploy Guide (`rewardchaind`)

This document explains how to **build the `rewardchaind` binary** for this Cosmos SDK chain and how to **deploy/run a node** (single-node validator or local multi-validator testnet).

## Modules Overview

This chain includes the following modules:

### Cosmos SDK Core Modules

- **`auth`** - Account management, authentication, and transaction signing
- **`bank`** - Token transfers and balance management
- **`staking`** - Validator staking, delegation, and bonding
- **`slashing`** - Validator slashing for downtime/double-signing
- **`distribution`** - Fee distribution and staking rewards
- **`mint`** - Token minting and inflation
- **`gov`** - On-chain governance (proposals and voting)
- **`params`** - Parameter management for modules
- **`upgrade`** - Chain upgrades and migrations
- **`crisis`** - Crisis invariant checks
- **`evidence`** - Evidence handling for misbehavior
- **`genutil`** - Genesis utilities and initialization
- **`consensus`** - Consensus parameter management

### Cosmos SDK Extended Modules

- **`authz`** - Authorization for granting permissions to other accounts
- **`feegrant`** - Fee grants (allowing one account to pay fees for another)
- **`group`** - Group-based governance and multisig
- **`nft`** - Non-fungible token (NFT) support
- **`vesting`** - Token vesting accounts
- **`circuit`** - Circuit breaker for emergency halts

### IBC Modules (Inter-Blockchain Communication)

- **`capability`** - Capability-based security for IBC
- **`ibc`** (core) - IBC core protocol
- **`transfer`** - IBC token transfers between chains
- **`interchain-accounts`** (ICA) - Cross-chain account management
- **`ibc-fee`** - Fee middleware for IBC transactions

### Custom Modules

- **`rewardchain`** - Partner management (create, disable, update partners with admin permissions)

### Module Execution Order

Modules execute in specific orders during block processing:

- **PreBlockers**: `upgrade` (handles upgrades before other modules)
- **BeginBlockers**: `mint`, `distribution`, `slashing`, `evidence`, `staking`, `authz`, `genutil`, `capability`, `ibc`, `transfer`, `ica`, `ibc-fee`, `rewardchain`
- **EndBlockers**: `crisis`, `gov`, `staking`, `feegrant`, `group`, `genutil`, `ibc`, `transfer`, `capability`, `ica`, `ibc-fee`, `rewardchain`

## Prerequisites

- **Go**: this repo targets **Go `1.24.x`** (see `go.mod`).
- **GNU Make** (optional but recommended).
- **git**

For Linux servers you'll also typically want:

- `build-essential` (or equivalent) for compiling
- `jq` (optional, helpful for config/genesis inspection)

## Token Economics & Chain Configuration

### Token Denomination

**Default Bond Denom**: `stake`

The chain uses `stake` as the default bond denomination (the token used for staking, delegation, and validator operations). This is the Cosmos SDK default and can be customized in your genesis configuration.

**To customize the denom**, you'll need to:

1. Set it in your genesis file's `staking` module params:
```json
{
  "app_state": {
    "staking": {
      "params": {
        "bond_denom": "your-token-denom"
      }
    }
  }
}
```

2. Update all references in genesis (balances, delegations, etc.) to use the new denom.

**Note**: The examples in this guide use `stake` as the denomination. Replace with your actual denom when deploying.

### Chain ID

**Default Chain ID**: `rewardchain`

The default chain ID is derived from the app name (`reward-chain`) with hyphens removed. For production deployments, you should **always explicitly set** a unique chain ID:

```bash
rewardchaind init "$MONIKER" --chain-id "rewardchain-1" --home "$HOME"
```

**Chain ID naming conventions**:
- Use a descriptive name (e.g., `rewardchain-mainnet`, `rewardchain-testnet`)
- Include network identifier (e.g., `rewardchain-1`, `rewardchain-2`)
- Keep it consistent across all nodes in the network

### Vesting Logic

The chain includes the **`vesting`** module, which supports token vesting accounts. Vesting accounts lock tokens for a specified period, releasing them gradually or all at once at the end of the vesting period.

**Supported vesting account types**:
- **Base Vesting Account**: Basic vesting with start/end times
- **Continuous Vesting Account**: Linear vesting over time
- **Delayed Vesting Account**: All tokens vest at end time
- **Periodic Vesting Account**: Multiple vesting periods

**Creating vesting accounts in genesis**:

You can create vesting accounts in your genesis file. The structure includes:

```json
{
  "app_state": {
    "auth": {
      "accounts": [
        {
          "@type": "/cosmos.vesting.v1beta1.ContinuousVestingAccount",
          "base_vesting_account": {
            "base_account": {
              "address": "reward1...",
              "pub_key": null,
              "account_number": "0",
              "sequence": "0"
            },
            "original_vesting": [
              {
                "denom": "stake",
                "amount": "1000000000"
              }
            ],
            "delegated_free": [],
            "delegated_vesting": [],
            "end_time": "1735689600"
          },
          "start_time": "1704153600"
        }
      ]
    }
  }
}
```

**Vesting parameters**:
- `original_vesting`: Total tokens to be vested
- `start_time`: Vesting start time (UNIX timestamp)
- `end_time`: Vesting end time (UNIX timestamp)
- `delegated_free`: Tokens that can be delegated immediately
- `delegated_vesting`: Tokens that are vesting and delegated

**Note**: Specific vesting schedules and parameters should be configured based on your tokenomics model. Please provide details on your vesting requirements for more specific guidance.

## Build the binary

### Build + install into your `GOBIN`

From the repo root:

```bash
cd /path/to/reward-chain
make install
```

This runs `go install ... ./cmd/reward-chaind` with version/commit ldflags and installs the binary as:

- `rewardchaind`

Verify:

```bash
rewardchaind version
which rewardchaind
```

If `rewardchaind` is installed but your shell can’t find it, add your Go bin dir to PATH (common default is `$(go env GOPATH)/bin`):

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

### Build a local binary into `./build/` (no install)

```bash
cd /path/to/reward-chain
mkdir -p build
go build -o ./build/rewardchaind ./cmd/reward-chaind
./build/rewardchaind version
```

### Reproducible “release-style” build notes

The Makefile injects version metadata via ldflags. If you need a fixed version string:

```bash
VERSION=v0.1.0 make install
rewardchaind version
```

## Single-node deployment (validator)

This is the common “one-machine, one-validator” flow. It produces a new chain from scratch.

### 1) Choose basics

Pick:

- **CHAIN_ID**: e.g. `rewardchain-1`
- **MONIKER**: e.g. `validator-1`
- **HOME**: node home directory, e.g. `~/.rewardchain`

> Note: by default, this chain sets the default `--chain-id` to `rewardchain` (hyphens removed from the app name). For real deployments you should **always explicitly set** `--chain-id`.

### 2) Initialize node files

```bash
export CHAIN_ID="rewardchain-1"
export MONIKER="validator-1"
export HOME="$HOME/.rewardchain"

rewardchaind init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME"
```

### 3) Create (or recover) a key

For a brand-new key:

```bash
rewardchaind keys add validator --home "$HOME" --keyring-backend os
```

To recover from an existing mnemonic:

```bash
rewardchaind keys add validator --recover --home "$HOME" --keyring-backend os
```

Get the account address:

```bash
VAL_ADDR="$(rewardchaind keys show validator -a --home "$HOME" --keyring-backend os)"
echo "$VAL_ADDR"
```

You should see addresses starting with:

- `reward1...` (account)
- `rewardvaloper1...` (validator operator)
- `rewardvalcons1...` (consensus)

### 4) Fund your account in genesis

Add a genesis balance for the validator account.

> Replace denominations/amounts as needed for your chain economics. `stake` is the Cosmos SDK default bond denom; your chain may choose a different denom later.

```bash
rewardchaind genesis add-genesis-account "$VAL_ADDR" 100000000stake --home "$HOME"
```

### 5) Create a validator gentx

```bash
rewardchaind genesis gentx validator 1000000stake \
  --chain-id "$CHAIN_ID" \
  --home "$HOME" \
  --keyring-backend os
```

### 6) Collect gentxs

```bash
rewardchaind genesis collect-gentxs --home "$HOME"
```

Optional sanity checks:

```bash
rewardchaind genesis validate-genesis --home "$HOME"
```

### 7) Configure networking (minimum)

Edit:

- `"$HOME/config/config.toml"`: P2P/RPC
- `"$HOME/config/app.toml"`: API/GRPC/min gas prices

At minimum, ensure `minimum-gas-prices` is non-empty (otherwise nodes can refuse to start depending on config).

You can also set it via flag at start time:

```bash
rewardchaind start --home "$HOME" --minimum-gas-prices "0.0001stake"
```

### 8) Start the node

```bash
rewardchaind start --home "$HOME"
```

## Linux server deployment with systemd

Below is a minimal template. Adjust user, paths, ports, pruning, and gas prices for your needs.

### 1) Put the binary on the server

Either build on-server using `make install`, or copy the built binary:

```bash
scp ./build/rewardchaind user@server:/usr/local/bin/rewardchaind
```

### 2) Create the systemd unit

Create `/etc/systemd/system/rewardchaind.service`:

```ini
[Unit]
Description=rewardchaind
After=network-online.target
Wants=network-online.target

[Service]
User=rewardchain
ExecStart=/usr/local/bin/rewardchaind start --home /home/rewardchain/.rewardchain --minimum-gas-prices 0.0001stake
Restart=on-failure
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

Enable + start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable rewardchaind
sudo systemctl start rewardchaind
sudo journalctl -u rewardchaind -f
```

## Local multi-validator network (`multi-node`)

This repo includes a built-in generator that creates multiple validator homes with keys, configs, a shared genesis, and gentxs.

### Generate N validator directories

```bash
rewardchaind multi-node --v 4 --output-dir ./.testnets
```

This will create directories like:

- `./.testnets/validator0`
- `./.testnets/validator1`
- `./.testnets/validator2`
- `./.testnets/validator3`

Each contains `config/` with `genesis.json`, `config.toml`, and `app.toml`.

### Start each node (separate terminals)

Example for validator 0:

```bash
rewardchaind start --home ./.testnets/validator0
```

Repeat for validator1..validatorN.

Notes:

- RPC ports are automatically spaced; see `config/config.toml` per node to confirm.
- Prometheus instrumentation is enabled by the generator.

## “In-place testnet” (advanced)

There is also an `in-place-testnet` command intended for **local testing from an existing node state** by replacing the validator set and funding some accounts.

Explore help:

```bash
rewardchaind in-place-testnet --help
```

## Troubleshooting

### “binary not found”

If you used `make install`, ensure your Go bin is on PATH:

```bash
go env GOPATH
go env GOBIN
```

Common default is `$(go env GOPATH)/bin`.

### “minimum gas price is not set”

Set `minimum-gas-prices` in `app.toml`, or start with:

```bash
rewardchaind start --minimum-gas-prices "0.0001stake"
```

### “wrong chain-id”

Always pass `--chain-id` for tx/gentx steps. If you created genesis with one chain-id and later used another, re-init your home directory and regenerate genesis.

### Bech32 prefix changes

If you changed the Bech32 prefix (e.g. `cosmos1...` → `reward1...`), you **must regenerate**:

- keys (or at least re-derive addresses under the new HRP)
- genesis accounts / gentxs / genesis file

Old Bech32 addresses won’t validate under the new prefix.

## Partner module usage

This chain includes a `rewardchain` module that manages on-chain partners.

### Admin addresses

Only **admin wallets** can create, disable, or update partners. Admin addresses are stored in module parameters and can be set:

#### 1) Set admin addresses in genesis (initial setup)

Edit your genesis file (`$HOME/config/genesis.json`) and add admin addresses to the `rewardchain` module params:

```json
{
  "app_state": {
    "rewardchain": {
      "params": {
        "admin_addresses": [
          "reward1abc...",
          "reward1def..."
        ]
      },
      "partner_list": []
    }
  }
}
```

Then validate and start your chain:

```bash
rewardchaind genesis validate-genesis --home "$HOME"
rewardchaind start --home "$HOME"
```

#### 2) Update admin addresses via governance

Since `UpdateParams` is controlled by the governance module, you must submit a governance proposal to change admin addresses.

**Step 1**: Create a proposal JSON file (`proposal-update-admin-addresses.json`):

```json
{
  "messages": [
    {
      "@type": "/rewardchain.rewardchain.MsgUpdateParams",
      "authority": "reward1...",  // governance module address (see note below)
      "params": {
        "admin_addresses": [
          "reward1abc123...",
          "reward1def456...",
          "reward1ghi789..."
        ]
      }
    }
  ],
  "metadata": "ipfs://CID",  // optional: IPFS hash of proposal metadata
  "deposit": "1000000stake"   // minimum deposit required (check your chain's min deposit)
}
```

**Note on authority**: The `authority` field should be the governance module account address. By default, this is derived from the `gov` module name. You can find it by:

```bash
# Query module accounts to find the gov module address
rewardchaind query auth module-accounts

# Or check your app configuration
# The default is typically: reward1... (derived from module name "gov")
```

**Step 2**: Submit the proposal:

```bash
rewardchaind tx gov submit-proposal proposal-update-admin-addresses.json \
  --from validator \
  --keyring-backend os \
  --chain-id "$CHAIN_ID" \
  --home "$HOME"
```

This will output a proposal ID (e.g., `1`).

**Step 3**: Vote on the proposal (if you have voting power):

```bash
PROPOSAL_ID=1  # Use the ID from step 2
rewardchaind tx gov vote $PROPOSAL_ID yes \
  --from validator \
  --keyring-backend os \
  --chain-id "$CHAIN_ID" \
  --home "$HOME"
```

**Step 4**: Wait for the voting period to end. Once the proposal passes, the admin addresses will be updated automatically.

**Query current params** (including admin addresses):

```bash
rewardchaind query rewardchain params
```

**Example output**:
```json
{
  "params": {
    "admin_addresses": [
      "reward1abc123...",
      "reward1def456..."
    ]
  }
}
```

**Verify admin access**: Only addresses in the `admin_addresses` list can execute partner mutations. If a non-admin tries to create/disable/update a partner, they will receive an `unauthorized` error.

### CLI (tx)

Create a partner:

```bash
rewardchaind tx rewardchain create-partner "Acme Inc" "retail" "Mumbai" "IN" \
  --from validator \
  --keyring-backend os \
  --chain-id "$CHAIN_ID" \
  --home "$HOME"
```

Disable a partner:

```bash
rewardchaind tx rewardchain disable-partner 1 \
  --from validator \
  --keyring-backend os \
  --chain-id "$CHAIN_ID" \
  --home "$HOME"
```

Update a partner:

```bash
rewardchaind tx rewardchain update-partner 1 "Acme Inc" "retail" "Mumbai" "IN" \
  --from validator \
  --keyring-backend os \
  --chain-id "$CHAIN_ID" \
  --home "$HOME"
```

### CLI (query)

List partners:

```bash
rewardchaind query rewardchain partners
```

Get partner by id:

```bash
rewardchaind query rewardchain partner 1
```

### HTTP API (gRPC-gateway)

When the node’s API server is enabled, these endpoints are available:

- `GET /reward-chain/rewardchain/partners`
- `GET /reward-chain/rewardchain/partners/{id}`



