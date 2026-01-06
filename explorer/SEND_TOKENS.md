# Send Tokens from Alice to Bob

This guide shows how to send tokens from the `alice` account to the `bob` account on the Reward Chain.

## Prerequisites

- Your chain must be running
- The `alice` and `bob` accounts must exist in your keyring
- You need to know your chain ID (default: `rewardchain`)

## Step 1: Get Account Addresses

First, get the addresses for both accounts:

```bash
# Get Alice's address
rewardchaind keys show alice --address

# Get Bob's address
rewardchaind keys show bob --address
```

Or get both addresses in one command:

```bash
# Get Alice's address
ALICE_ADDRESS=$(rewardchaind keys show alice --address)

# Get Bob's address
BOB_ADDRESS=$(rewardchaind keys show bob --address)

# Display them
echo "Alice: $ALICE_ADDRESS"
echo "Bob: $BOB_ADDRESS"
```

## Step 2: Check Balances (Optional)

Before sending, check the current balances:

```bash
# Check Alice's balance
rewardchaind query bank balances $(rewardchaind keys show alice --address)

# Check Bob's balance
rewardchaind query bank balances $(rewardchaind keys show bob --address)
```

## Step 3: Send Tokens

### Basic Command (Test Keyring)

If you're using the `test` keyring backend (default for local development):

```bash
rewardchaind tx bank send alice $(rewardchaind keys show bob --address) 1000stake \
  --chain-id rewardchain \
  --keyring-backend test \
  --from alice \
  --yes
```

### With OS Keyring Backend

If you're using the `os` keyring backend:

```bash
rewardchaind tx bank send alice $(rewardchaind keys show bob --address) 1000stake \
  --chain-id rewardchain \
  --keyring-backend os \
  --from alice \
  --yes
```

### Send Multiple Token Types

To send both `stake` and `token`:

```bash
rewardchaind tx bank send alice $(rewardchaind keys show bob --address) 1000stake,500token \
  --chain-id rewardchain \
  --keyring-backend test \
  --from alice \
  --yes
```

### With Gas Settings

To specify gas price and gas limit:

```bash
rewardchaind tx bank send alice $(rewardchaind keys show bob --address) 1000stake \
  --chain-id rewardchain \
  --keyring-backend test \
  --from alice \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.0001stake \
  --yes
```

### With Memo

To add a memo to the transaction:

```bash
rewardchaind tx bank send alice $(rewardchaind keys show bob --address) 1000stake \
  --chain-id rewardchain \
  --keyring-backend test \
  --from alice \
  --memo "Payment from Alice to Bob" \
  --yes
```

## Complete Example Script

Here's a complete script that does everything:

```bash
#!/bin/bash

# Set variables
CHAIN_ID="rewardchain"
KEYRING_BACKEND="test"  # Change to "os" if using OS keyring
AMOUNT="1000stake"

# Get addresses
ALICE_ADDRESS=$(rewardchaind keys show alice --address --keyring-backend $KEYRING_BACKEND)
BOB_ADDRESS=$(rewardchaind keys show bob --address --keyring-backend $KEYRING_BACKEND)

echo "Alice address: $ALICE_ADDRESS"
echo "Bob address: $BOB_ADDRESS"
echo ""

# Check balances before
echo "=== Balances Before ==="
echo "Alice:"
rewardchaind query bank balances $ALICE_ADDRESS --chain-id $CHAIN_ID
echo ""
echo "Bob:"
rewardchaind query bank balances $BOB_ADDRESS --chain-id $CHAIN_ID
echo ""

# Send tokens
echo "=== Sending $AMOUNT from Alice to Bob ==="
rewardchaind tx bank send alice $BOB_ADDRESS $AMOUNT \
  --chain-id $CHAIN_ID \
  --keyring-backend $KEYRING_BACKEND \
  --from alice \
  --yes

# Wait for transaction to be included
sleep 3

# Check balances after
echo ""
echo "=== Balances After ==="
echo "Alice:"
rewardchaind query bank balances $ALICE_ADDRESS --chain-id $CHAIN_ID
echo ""
echo "Bob:"
rewardchaind query bank balances $BOB_ADDRESS --chain-id $CHAIN_ID
```

## One-Liner Command

Quick one-liner to send tokens:

```bash
rewardchaind tx bank send alice $(rewardchaind keys show bob --address) 1000stake --chain-id rewardchain --keyring-backend test --from alice --yes
```

## Verify Transaction

After sending, you can verify the transaction:

```bash
# Get the transaction hash from the output, then query it
rewardchaind query tx <TX_HASH> --chain-id rewardchain
```

Or check the latest transactions:

```bash
# Query recent transactions for Alice
rewardchaind query txs --events 'message.sender=reward1...' --chain-id rewardchain --limit 10
```

## Troubleshooting

### Error: "account sequence mismatch"

This means the account sequence is out of sync. Wait a moment and try again, or check the account sequence:

```bash
rewardchaind query account $(rewardchaind keys show alice --address) --chain-id rewardchain
```

### Error: "insufficient funds"

Check the account balance:

```bash
rewardchaind query bank balances $(rewardchaind keys show alice --address) --chain-id rewardchain
```

### Error: "key not found"

Make sure the key exists in your keyring:

```bash
# List all keys
rewardchaind keys list --keyring-backend test
```

### Error: "chain-id mismatch"

Make sure you're using the correct chain ID. Check your chain's ID:

```bash
rewardchaind status | jq .node_info.network
```

## Token Denominations

Based on the chain configuration:
- **stake**: Base staking token (6 decimals: 1 STAKE = 1,000,000 stake)
- **token**: Additional token (6 decimals: 1 TOKEN = 1,000,000 token)

Examples:
- `1000stake` = 1000 base units
- `1000000stake` = 1 STAKE (display unit)
- `500token` = 500 base units

## Additional Resources

- [Cosmos SDK Bank Module Docs](https://docs.cosmos.network/main/modules/bank)
- [Reward Chain Documentation](../docs/)

