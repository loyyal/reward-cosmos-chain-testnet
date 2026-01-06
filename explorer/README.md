# Ping.pub Explorer Setup for Reward Chain

This directory contains the configuration and setup scripts for running a [Ping.pub](https://ping.pub) block explorer for the Reward Chain.

## Overview

Ping.pub is a lightweight, easy-to-setup block explorer perfect for local development and testnets. It provides a web interface to browse blocks, transactions, validators, and accounts on your chain.

## Prerequisites

Before setting up the explorer, ensure you have:

1. **Node.js** (v18 or later) - [Install Node.js](https://nodejs.org/)
2. **Yarn** or **npm** - Package manager (Yarn is preferred)
3. **Git** - For cloning the explorer repository
4. **Running Reward Chain node** - Your chain must be running with:
   - RPC endpoint: `http://localhost:26657`
   - REST API endpoint: `http://localhost:1317`

## Quick Start

### Automatic Setup

The easiest way to set up the explorer is using the provided setup script:

```bash
cd explorer
chmod +x setup.sh
./setup.sh setup
```

This script will:
1. Check for required dependencies (Node.js, Yarn/npm, Git)
2. Clone the Ping.pub explorer repository
3. Configure it for Reward Chain
4. Install dependencies
5. Start the development server

The explorer will be available at: **http://localhost:5173/reward-chain**

### Manual Setup

If you prefer to set up manually or the automatic script fails:

```bash
# 1. Clone the Ping.pub explorer
git clone https://github.com/ping-pub/explorer.git ~/ping-pub-explorer
cd ~/ping-pub-explorer

# 2. Install dependencies
yarn install  # or npm install --legacy-peer-deps
# Note: npm may require --legacy-peer-deps flag to resolve dependency conflicts

# 3. Create chain configuration directory
mkdir -p chains/mainnet

# 4. Copy the Reward Chain configuration
cp /path/to/reward-chain/explorer/ping_pub_config.json chains/mainnet/reward-chain.json

# 5. Start the development server
yarn dev  # or npm run dev
```

## Configuration

The chain configuration is stored in `ping_pub_config.json`. Here's what each field means:

### Chain Information

- **chain_name**: `reward-chain` - Display name for the chain
- **registry_name**: `reward-chain` - Registry identifier
- **sdk_version**: `0.50.11` - Cosmos SDK version used by the chain
- **coin_type**: `118` - HD wallet coin type (standard for Cosmos chains)

### Network Endpoints

- **api**: REST API endpoints (default: `http://localhost:1317`)
- **rpc**: Tendermint RPC endpoints (default: `http://localhost:26657`)

### Address Configuration

- **addr_prefix**: `reward` - Address prefix for accounts
- **bech32_prefix**: `reward` - Bech32 encoding prefix

### Token Configuration

- **base**: `stake` - Base denomination (smallest unit)
- **symbol**: `STAKE` - Display symbol
- **exponent**: `6` - Decimal places (1 STAKE = 1,000,000 stake)

### Features

- **ibc-transfer**: Inter-Blockchain Communication token transfers
- **ibc-go**: IBC Go module support
- **keplr_features**: Keplr wallet integration features

## CORS Configuration

For the explorer to connect to your local node, you need to enable CORS in your node's configuration.

### Enable CORS in config.toml

Edit `~/.rewardchain/config/config.toml`:

```toml
[rpc]
cors_allowed_origins = ["*"]
```

### Enable CORS in app.toml

Edit `~/.rewardchain/config/app.toml`:

```toml
[api]
enable = true
enabled-unsafe-cors = true
```

After making these changes, **restart your node**.

### Quick CORS Help

Run the setup script with the `cors` option to see CORS configuration instructions:

```bash
./setup.sh cors
```

## Usage

### Starting the Explorer

```bash
# Using the setup script
./setup.sh setup

# Or manually
cd ~/ping-pub-explorer
yarn dev  # or npm run dev
```

### Accessing the Explorer

Once running, open your browser and navigate to:

- **Local**: http://localhost:5173/reward-chain
- **Network**: http://YOUR_IP:5173/reward-chain (if accessible on your network)

### Explorer Features

The Ping.pub explorer provides:

- **Blocks**: Browse recent blocks and block details
- **Transactions**: View transaction history and details
- **Validators**: See validator information and staking details
- **Accounts**: Check account balances and transaction history
- **Proposals**: View governance proposals (if applicable)
- **IBC**: Monitor IBC transfers and channels

## Troubleshooting

### Port Already in Use

If port 5173 is already in use:

```bash
# Find what's using the port
lsof -i :5173

# Kill the process or change the port in the explorer's configuration
```

### Dependency Resolution Errors (npm)

If you encounter `ERESOLVE unable to resolve dependency tree` errors when using npm:

```bash
# Use --legacy-peer-deps flag
npm install --legacy-peer-deps

# Or use yarn instead (recommended)
yarn install
```

This is a known issue with the Ping.pub explorer's dependencies. The `--legacy-peer-deps` flag tells npm to use the legacy peer dependency resolution algorithm, which is more permissive.

### CORS Errors

If you see CORS errors in the browser console:

1. Verify CORS is enabled in both `config.toml` and `app.toml`
2. Restart your node after making changes
3. Check that the API server is enabled: `enable = true` in `app.toml`

### Node Not Responding

If the explorer can't connect to your node:

1. Verify your node is running:
   ```bash
   curl http://localhost:26657/status
   curl http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info
   ```

2. Check firewall settings if accessing remotely

3. Verify the endpoints in `ping_pub_config.json` match your node's configuration

### Explorer Shows "Chain Not Found"

1. Verify the configuration file is in the correct location:
   ```bash
   ls ~/ping-pub-explorer/chains/mainnet/reward-chain.json
   ```

2. Check that the JSON is valid:
   ```bash
   cat ~/ping-pub-explorer/chains/mainnet/reward-chain.json | jq .
   ```

3. Restart the explorer after adding/updating the configuration

## Production Deployment

For production deployments, you'll want to:

1. **Build for production**:
   ```bash
   cd ~/ping-pub-explorer
   yarn build  # or npm run build
   ```

2. **Use a production web server** (nginx, Apache, etc.) to serve the built files

3. **Update endpoints** in `ping_pub_config.json` to use production RPC/API URLs

4. **Configure HTTPS** for secure access

5. **Set up proper CORS** origins instead of `["*"]` for security

## Configuration for Different Networks

To use the explorer with different networks (testnet, mainnet, etc.):

1. Create separate configuration files:
   - `chains/testnet/reward-chain-testnet.json`
   - `chains/mainnet/reward-chain-mainnet.json`

2. Update the endpoints in each configuration file

3. Access via:
   - `http://localhost:5173/reward-chain-testnet`
   - `http://localhost:5173/reward-chain-mainnet`

## Additional Resources

- [Ping.pub Explorer GitHub](https://github.com/ping-pub/explorer)
- [Ping.pub Website](https://ping.pub)
- [Cosmos SDK Documentation](https://docs.cosmos.network)
- [Reward Chain Documentation](../docs/)

## Support

For issues specific to:
- **Ping.pub Explorer**: Check the [Ping.pub GitHub Issues](https://github.com/ping-pub/explorer/issues)
- **Reward Chain**: Check the main repository documentation

## Script Options

The `setup.sh` script supports several options:

```bash
./setup.sh setup   # Full automatic setup and start
./setup.sh manual  # Print manual setup instructions
./setup.sh cors    # Print CORS configuration instructions
```

## License

This configuration is part of the Reward Chain project. The Ping.pub explorer is maintained by the Ping.pub team.

