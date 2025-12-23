# Block Explorer Setup Guide

This guide explains how to set up a block explorer for your local `rewardchain` chain.

## Prerequisites

- Your local chain must be running (see `build-and-deploy.md`)
- Docker and Docker Compose installed (for Docker-based explorers)
- Node.js 18+ (for some explorer options)

## Chain Endpoints

Your chain exposes the following endpoints by default:

- **RPC**: `http://localhost:26657` (Tendermint RPC)
- **REST API**: `http://localhost:1317` (Cosmos SDK REST API)
- **gRPC**: `localhost:9090` (Cosmos SDK gRPC)

Verify these are accessible:

```bash
# Check RPC
curl http://localhost:26657/status

# Check REST API
curl http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info
```

## Option 1: Ping.pub Explorer (Recommended for Local Development)

Ping.pub is a lightweight, easy-to-setup explorer perfect for local chains.

### Setup Steps

1. **Clone the Ping.pub repository**:

```bash
git clone https://github.com/ping-pub/explorer.git
cd explorer
```

2. **Install dependencies**:

```bash
npm install
```

3. **Configure for your chain**:

Create or edit `src/chains/rewardchain.ts`:

```typescript
export default {
  chainId: 'rewardchain',
  chainName: 'Reward Chain',
  addressPrefix: 'reward',
  rpc: 'http://localhost:26657',
  rest: 'http://localhost:1317',
  coinType: 118,
  coinDenom: 'stake',
  coinMinimalDenom: 'stake',
  coinDecimals: 6,
  bech32Config: {
    bech32PrefixAccAddr: 'reward',
    bech32PrefixAccPub: 'rewardpub',
    bech32PrefixValAddr: 'rewardvaloper',
    bech32PrefixValPub: 'rewardvaloperpub',
    bech32PrefixConsAddr: 'rewardvalcons',
    bech32PrefixConsPub: 'rewardvalconspub',
  },
}
```

4. **Add your chain to the chains list**:

Edit `src/chains/index.ts` and add:

```typescript
import rewardchain from './rewardchain'
// ... other imports

export const chains = [
  // ... other chains
  rewardchain,
]
```

5. **Start the explorer**:

```bash
npm run dev
```

The explorer will be available at `http://localhost:3000`

## Option 2: BlockScout (Full-Featured Explorer)

BlockScout is a more comprehensive explorer with advanced features.

### Setup with Docker Compose

1. **Create a `docker-compose.blockscout.yml` file**:

```yaml
version: '3.8'

services:
  blockscout:
    image: blockscout/blockscout:latest
    ports:
      - "4000:4000"
    environment:
      - ETHEREUM_JSONRPC_VARIANT=geth
      - ETHEREUM_JSONRPC_HTTP_URL=http://localhost:26657
      - ETHEREUM_JSONRPC_WS_URL=ws://localhost:26657/websocket
      - DATABASE_URL=postgresql://blockscout:blockscout@postgres:5432/blockscout
      - SECRET_KEY_BASE=your_secret_key_here
      - CHAIN_ID=rewardchain
      - COIN=stake
      - NETWORK=Reward Chain
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_USER=blockscout
      - POSTGRES_PASSWORD=blockscout
      - POSTGRES_DB=blockscout
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine

volumes:
  postgres_data:
```

**Note**: BlockScout is primarily designed for EVM chains. For Cosmos SDK chains, you may need to use a Cosmos-specific explorer instead.

## Option 3: Cosmos Explorer (Big Dipper / Mintscan-style)

For Cosmos SDK chains, you can use a Cosmos-specific explorer. Here are two options:

### Option 3a: Big Dipper

Big Dipper is a popular Cosmos explorer. However, it requires more setup and is typically used for production networks.

### Option 3b: Simple Custom Explorer

Create a simple explorer using the REST API:

1. **Create a basic HTML explorer** (`explorer/index.html`):

```html
<!DOCTYPE html>
<html>
<head>
    <title>Reward Chain Explorer</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .block { border: 1px solid #ddd; padding: 10px; margin: 10px 0; }
        .tx { background: #f5f5f5; padding: 5px; margin: 5px 0; }
    </style>
</head>
<body>
    <h1>Reward Chain Explorer</h1>
    <div id="status"></div>
    <div id="blocks"></div>
    
    <script>
        const REST_API = 'http://localhost:1317';
        const RPC_API = 'http://localhost:26657';
        
        async function fetchStatus() {
            try {
                const response = await fetch(`${RPC_API}/status`);
                const data = await response.json();
                document.getElementById('status').innerHTML = `
                    <h2>Chain Status</h2>
                    <p>Chain ID: ${data.result.node_info.network}</p>
                    <p>Latest Block: ${data.result.sync_info.latest_block_height}</p>
                `;
            } catch (error) {
                console.error('Error fetching status:', error);
            }
        }
        
        async function fetchLatestBlocks() {
            try {
                const response = await fetch(`${REST_API}/cosmos/base/tendermint/v1beta1/blocks/latest`);
                const data = await response.json();
                document.getElementById('blocks').innerHTML = `
                    <h2>Latest Block</h2>
                    <div class="block">
                        <p>Height: ${data.block.header.height}</p>
                        <p>Time: ${new Date(data.block.header.time).toLocaleString()}</p>
                        <p>Hash: ${data.block_id.hash}</p>
                    </div>
                `;
            } catch (error) {
                console.error('Error fetching blocks:', error);
            }
        }
        
        // Fetch data on load and every 5 seconds
        fetchStatus();
        fetchLatestBlocks();
        setInterval(() => {
            fetchStatus();
            fetchLatestBlocks();
        }, 5000);
    </script>
</body>
</html>
```

2. **Serve it with a simple HTTP server**:

```bash
# Using Python
python3 -m http.server 8080

# Or using Node.js
npx http-server -p 8080
```

3. **Open in browser**: `http://localhost:8080`

## Option 4: Using Ignite CLI Explorer (If Available)

If you're using Ignite CLI, it may have built-in explorer support. Check if there's an explorer command:

```bash
ignite chain serve --explorer
```

Or check the Ignite documentation for explorer integration.

## Recommended Approach for Local Development

For quick local development, I recommend:

1. **Start with Option 3b** (Simple Custom Explorer) for immediate visibility
2. **Move to Option 1** (Ping.pub) for a more polished experience
3. **Use Option 2 or 3a** for production-like testing

## Troubleshooting

### CORS Issues

If you encounter CORS errors when accessing the REST API from a browser, you may need to enable CORS in your node's configuration.

Edit `~/.rewardchain/config/app.toml` and ensure:

```toml
[api]
# Enable CORS
enabled-unsafe-cors = true
```

Then restart your node.

### Port Conflicts

If ports are already in use:

- Check which ports your node is using: `rewardchaind config show`
- Adjust explorer configuration to use different ports
- Or change node ports in `config.toml` and `app.toml`

### API Not Responding

Ensure your node's API server is enabled:

```bash
# Check app.toml
cat ~/.rewardchain/config/app.toml | grep -A 5 "\[api\]"
```

Should show:
```toml
[api]
enable = true
address = "tcp://0.0.0.0:1317"
```

## Next Steps

- Customize the explorer UI to match your chain's branding
- Add more features (transaction search, address lookup, etc.)
- Set up monitoring and alerts
- Consider deploying a public explorer for testnet/mainnet

