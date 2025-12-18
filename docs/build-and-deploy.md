# Build & Deploy Guide (`rewardchaind`)

This document explains how to **build the `rewardchaind` binary** for this Cosmos SDK chain and how to **deploy/run a node** (single-node validator or local multi-validator testnet).

## Prerequisites

- **Go**: this repo targets **Go `1.24.x`** (see `go.mod`).
- **GNU Make** (optional but recommended).
- **git**

For Linux servers you’ll also typically want:

- `build-essential` (or equivalent) for compiling
- `jq` (optional, helpful for config/genesis inspection)

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


