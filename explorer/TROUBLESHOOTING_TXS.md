# Troubleshooting: Transactions Not Showing in Ping.pub Explorer

If ping.pub explorer is working but transactions are not appearing, follow these steps to diagnose and fix the issue.

## Quick Diagnosis

Run these commands to quickly check the common issues:

```bash
# 1. Check if API is enabled
grep "enable = true" ~/.rewardchain/config/app.toml | grep -v "^#"

# 2. Check if CORS is enabled
grep "enabled-unsafe-cors = true" ~/.rewardchain/config/app.toml | grep -v "^#"

# 3. Check transaction indexing
grep 'indexer = "kv"' ~/.rewardchain/config/config.toml | grep -v "^#"

# 4. Test API endpoint
curl -s http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info | head -20

# 5. Test transaction query
curl -s "http://localhost:1317/cosmos/tx/v1beta1/txs?pagination.limit=5" | head -30
```

## Step 1: Verify API Server is Enabled

The API server must be enabled for the explorer to fetch transaction data.

### Check Current API Configuration

```bash
cat ~/.rewardchain/config/app.toml | grep -A 10 "\[api\]"
```

Should show:
```toml
[api]
enable = true
address = "tcp://0.0.0.0:1317"
swagger = true
enabled-unsafe-cors = true
```

### Enable API Server (if not enabled)

Edit `~/.rewardchain/config/app.toml` and ensure:

```toml
[api]
enable = true
address = "tcp://0.0.0.0:1317"
swagger = true
enabled-unsafe-cors = true
max-open-connections = 1000
rpc-read-timeout = 10
rpc-write-timeout = 10
```

**⚠️ IMPORTANT: Restart your node after making changes!**

## Step 2: Verify CORS is Enabled

CORS must be enabled for the browser-based explorer to access the API.

### Check CORS Settings

```bash
# Check RPC CORS in config.toml
cat ~/.rewardchain/config/config.toml | grep -A 2 "cors_allowed_origins"

# Check API CORS in app.toml
cat ~/.rewardchain/config/app.toml | grep "enabled-unsafe-cors"
```

### Enable CORS

**In `~/.rewardchain/config/config.toml`:**
```toml
[rpc]
cors_allowed_origins = ["*"]
```

**In `~/.rewardchain/config/app.toml`:**
```toml
[api]
enabled-unsafe-cors = true
```

**⚠️ Restart your node after making changes!**

## Step 3: Enable Transaction Indexing

Transactions must be indexed for the explorer to find them.

### Check Current Indexing Configuration

```bash
cat ~/.rewardchain/config/config.toml | grep -A 5 "\[tx_index\]"
```

### Enable Transaction Indexing

In `~/.rewardchain/config/config.toml`:

```toml
[tx_index]
indexer = "kv"
```

Options:
- `"kv"` - Key-value indexer (default, recommended) ✅
- `"null"` - No indexing (transactions won't be searchable) ❌
- `"psql"` - PostgreSQL indexer (requires database setup)

**⚠️ If you change from "null" to "kv", you may need to re-index. Restart your node!**

## Step 4: Test API Endpoints

Verify the API endpoints are responding correctly:

```bash
# Test RPC endpoint
echo "=== RPC Status ==="
curl -s http://localhost:26657/status | jq .result.node_info.network

# Test API endpoint - node info
echo "=== API Node Info ==="
curl -s http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info | jq .default_node_info.network

# Test API endpoint - latest block
echo "=== Latest Block ==="
curl -s http://localhost:1317/cosmos/base/tendermint/v1beta1/blocks/latest | jq .block.header.height

# Test API endpoint - recent transactions
echo "=== Recent Transactions ==="
curl -s "http://localhost:1317/cosmos/tx/v1beta1/txs?pagination.limit=5" | jq '.txs | length'
```

If any of these fail, the API is not properly configured.

## Step 5: Verify Transactions Exist on Chain

Check if transactions are actually on-chain:

```bash
# Get latest block with transactions
rewardchaind query block | jq .block.data.txs

# Get recent transactions
rewardchaind query txs --events 'message.action=send' --limit 10

# Query specific account transactions
ALICE_ADDR=$(rewardchaind keys show alice --address --keyring-backend test)
rewardchaind query txs --events "message.sender=$ALICE_ADDR" --limit 10
```

## Step 6: Check ping_pub_config.json Endpoints

Ensure the endpoints in your config match your node:

```bash
cat explorer/ping_pub_config.json | jq '.api, .rpc'
```

Should show:
```json
{
  "api": [
    {
      "address": "http://localhost:1317",
      "provider": "local"
    }
  ],
  "rpc": [
    {
      "address": "http://localhost:26657",
      "provider": "local"
    }
  ]
}
```

If your node is running on different ports, update accordingly.

## Step 7: Check Browser Console

Open the browser developer console (F12) and check for errors:

1. **CORS errors**: If you see CORS errors, CORS is not properly enabled
2. **404 errors**: API endpoints might be wrong
3. **Network errors**: Node might not be running or accessible
4. **Empty responses**: Transactions might not be indexed

## Step 8: Complete Fix Script

Run this script to automatically fix common issues:

```bash
#!/bin/bash

CONFIG_DIR="$HOME/.rewardchain/config"
APP_TOML="$CONFIG_DIR/app.toml"
CONFIG_TOML="$CONFIG_DIR/config.toml"

echo "=== Fixing Ping.pub Transaction Issues ==="
echo ""

# Backup configs
cp "$APP_TOML" "$APP_TOML.backup"
cp "$CONFIG_TOML" "$CONFIG_TOML.backup"
echo "✓ Backed up config files"

# Fix API settings
echo ""
echo "=== Fixing API Settings ==="
if grep -q "\[api\]" "$APP_TOML"; then
    # Enable API if disabled
    sed -i.bak 's/^enable = false/enable = true/' "$APP_TOML"
    # Enable CORS
    sed -i.bak 's/^enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "$APP_TOML"
    if ! grep -q "enabled-unsafe-cors" "$APP_TOML"; then
        sed -i.bak '/\[api\]/a enabled-unsafe-cors = true' "$APP_TOML"
    fi
    echo "✓ API settings updated"
else
    echo "⚠️  [api] section not found, you may need to add it manually"
fi

# Fix CORS in config.toml
echo ""
echo "=== Fixing CORS in config.toml ==="
if grep -q "cors_allowed_origins" "$CONFIG_TOML"; then
    sed -i.bak 's/cors_allowed_origins = \[.*\]/cors_allowed_origins = ["*"]/' "$CONFIG_TOML"
    echo "✓ CORS updated in config.toml"
else
    # Add CORS if not present
    if grep -q "\[rpc\]" "$CONFIG_TOML"; then
        sed -i.bak '/\[rpc\]/a cors_allowed_origins = ["*"]' "$CONFIG_TOML"
        echo "✓ CORS added to config.toml"
    fi
fi

# Fix transaction indexing
echo ""
echo "=== Fixing Transaction Indexing ==="
if grep -q "\[tx_index\]" "$CONFIG_TOML"; then
    sed -i.bak 's/indexer = "null"/indexer = "kv"/' "$CONFIG_TOML"
    sed -i.bak 's/indexer = "psql"/indexer = "kv"/' "$CONFIG_TOML"
    echo "✓ Transaction indexing set to 'kv'"
else
    # Add tx_index section if not present
    echo "" >> "$CONFIG_TOML"
    echo "[tx_index]" >> "$CONFIG_TOML"
    echo 'indexer = "kv"' >> "$CONFIG_TOML"
    echo "✓ Transaction indexing section added"
fi

echo ""
echo "=== Configuration Updated ==="
echo ""
echo "⚠️  IMPORTANT: Restart your node for changes to take effect!"
echo ""
echo "To restart:"
echo "  1. Stop your current node (Ctrl+C)"
echo "  2. Start it again: rewardchaind start --home ~/.rewardchain"
echo ""
echo "Then restart ping.pub explorer and refresh your browser."
```

## Step 9: Manual Configuration Edits

If the script doesn't work, manually edit the files:

### Edit `~/.rewardchain/config/app.toml`

Find the `[api]` section and ensure:

```toml
[api]
enable = true
address = "tcp://0.0.0.0:1317"
swagger = true
enabled-unsafe-cors = true
```

### Edit `~/.rewardchain/config/config.toml`

Find the `[rpc]` section and ensure:

```toml
[rpc]
cors_allowed_origins = ["*"]
```

Find the `[tx_index]` section and ensure:

```toml
[tx_index]
indexer = "kv"
```

## Step 10: Verify After Restart

After restarting your node, verify everything works:

```bash
# 1. Check API is responding
curl -s http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info | jq .default_node_info.network

# 2. Check transactions are queryable
curl -s "http://localhost:1317/cosmos/tx/v1beta1/txs?pagination.limit=1" | jq '.txs | length'

# 3. Send a test transaction
rewardchaind tx bank send alice $(rewardchaind keys show bob --address --keyring-backend test) 1000stake \
  --chain-id rewardchain \
  --keyring-backend test \
  --from alice \
  --yes

# 4. Wait a few seconds, then check if it appears
sleep 5
curl -s "http://localhost:1317/cosmos/tx/v1beta1/txs?pagination.limit=5" | jq '.txs[0].txhash'
```

## Common Issues and Solutions

### Issue: "API endpoint not responding"
**Solution**: Check if API is enabled and node is running:
```bash
rewardchaind status
curl http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info
```

### Issue: "CORS errors in browser console"
**Solution**: Enable CORS in both `app.toml` and `config.toml`, then restart node.

### Issue: "Transactions exist but don't show in explorer"
**Solution**: Enable transaction indexing (`indexer = "kv"` in `config.toml`).

### Issue: "Explorer shows old transactions but not new ones"
**Solution**: 
1. Check if node is synced: `rewardchaind status`
2. Verify transaction indexing is enabled
3. Clear browser cache and refresh

### Issue: "Empty transaction list"
**Solution**: 
1. Verify transactions exist: `rewardchaind query txs --limit 10`
2. Check API endpoint: `curl http://localhost:1317/cosmos/tx/v1beta1/txs?pagination.limit=5`
3. Ensure transaction indexing is enabled

## Still Not Working?

If transactions still don't appear after following all steps:

1. **Check node logs** for errors:
   ```bash
   # If running in terminal, check for errors
   # Or check log files if using systemd
   journalctl -u rewardchaind -f
   ```

2. **Verify ping.pub explorer version** - older versions might have compatibility issues

3. **Try querying transactions directly** via API to isolate the issue:
   ```bash
   # Get a transaction hash from a recent tx
   TX_HASH="YOUR_TX_HASH"
   curl http://localhost:1317/cosmos/tx/v1beta1/txs/$TX_HASH
   ```

4. **Check if it's a ping.pub specific issue** by testing the API directly in browser:
   ```
   http://localhost:1317/cosmos/tx/v1beta1/txs?pagination.limit=10
   ```

## Additional Resources

- [Reward Chain Documentation](../docs/)
- [Ping.pub Explorer GitHub](https://github.com/ping-pub/explorer)
- [Cosmos SDK API Documentation](https://docs.cosmos.network/main/core/grpc_rest)

