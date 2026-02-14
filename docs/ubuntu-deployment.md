# Ubuntu Server Deployment Guide

This guide provides step-by-step instructions for deploying the `rewardchaind` binary to an Ubuntu server and running it as a production validator node.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Server Setup](#server-setup)
3. [Building the Binary](#building-the-binary)
4. [Node Initialization](#node-initialization)
5. [Configuration](#configuration)
6. [Systemd Service Setup](#systemd-service-setup)
7. [Firewall Configuration](#firewall-configuration)
8. [Block Explorer Setup](#block-explorer-setup)
9. [Monitoring and Maintenance](#monitoring-and-maintenance)
10. [Troubleshooting](#troubleshooting)

## Prerequisites

### Server Requirements

- **Operating System**: Ubuntu 20.04 LTS or later (22.04 LTS recommended)
- **CPU**: Minimum 2 cores (4+ cores recommended)
- **RAM**: Minimum 4GB (8GB+ recommended)
- **Storage**: Minimum 100GB SSD (500GB+ recommended for production)
- **Network**: Stable internet connection with static IP address (recommended)

### Required Software

Before deploying, ensure your Ubuntu server has the following installed:

```bash
# Update system packages
sudo apt-get update
sudo apt-get upgrade -y

# Install essential build tools
sudo apt-get install -y build-essential git curl wget jq

# Install Go (if building on server)
# Download Go 1.24.x (check https://go.dev/dl/ for latest version)
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version
```

## Server Setup

### 1. Set Up Directory Structure

This guide assumes you're running as the `ubuntu` user. If you prefer to use a dedicated user, you can create one, but all commands below will work with the `ubuntu` user.

```bash
# Create directories
mkdir -p ~/rewardchain/{bin,config,data,logs}
mkdir -p ~/.rewardchain

# Set proper permissions
chmod 700 ~/.rewardchain
```

## Building the Binary

You have two options: build on the server or build locally and transfer.

### Option 1: Build on the Server

If you have Go installed on the server:

```bash
# Clone the repository
cd ~
git clone https://github.com/your-username/rewardchain.git
cd rewardchain

# Build the binary
make install

# Verify the binary
rewardchaind version

# Copy binary to a standard location
sudo cp $(which rewardchaind) /usr/local/bin/rewardchaind
sudo chmod +x /usr/local/bin/rewardchaind
```

### Option 2: Build Locally and Transfer

Build the binary on your local machine:

```bash
# On your local machine
cd /path/to/rewardchain
make install

# Build a release binary
mkdir -p build
go build -o ./build/rewardchaind ./cmd/rewardchaind

# Verify
./build/rewardchaind version
```

Transfer to the server:

```bash
# From your local machine
scp ./build/rewardchaind ubuntu@your-server-ip:/tmp/rewardchaind

# On the server (as ubuntu user)
sudo mv /tmp/rewardchaind /usr/local/bin/rewardchaind
sudo chmod +x /usr/local/bin/rewardchaind
sudo chown root:root /usr/local/bin/rewardchaind

# Verify
rewardchaind version
```

## Node Initialization

### 1. Initialize the Node

```bash
# Set variables
export CHAIN_ID="rewardchain-1"  # Change to your chain ID
export MONIKER="validator-1"  # Your validator name
export REWARDCHAIN_HOME="$HOME/.rewardchain"

# DANGER: wipes the node’s local data & config
sudo rm -rf "$REWARDCHAIN_HOME"

mkdir -p "$REWARDCHAIN_HOME"

# Initialize the node
rewardchaind init "$MONIKER" --chain-id "$CHAIN_ID" --home "$REWARDCHAIN_HOME"
```

### 2. Create or Import Validator Key

**Option A: Create a new key**

```bash
rewardchaind keys add validator \
  --home "$REWARDCHAIN_HOME" \
  --keyring-backend file
```

**Important**: Save the mnemonic phrase in a secure location. You'll need it to recover your key.

**Option B: Import existing key from mnemonic**

```bash
rewardchaind keys add validator \
  --recover \
  --home "$REWARDCHAIN_HOME" \
  --keyring-backend file
```

**Get your validator address:**

```bash
VAL_ADDR=$(rewardchaind keys show validator -a --home "$REWARDCHAIN_HOME" --keyring-backend file)
echo "Validator address: $VAL_ADDR"
```

### 3. Configure Genesis (For New Chain)

If you're starting a new chain:

```bash
# Add genesis account
rewardchaind genesis add-genesis-account validator 1000000000stake --home "$REWARDCHAIN_HOME" --keyring-backend file

# Create validator gentx
# rewardchaind genesis gentx validator 1000000stake --chain-id "$CHAIN_ID" --home "$REWARDCHAIN_HOME" --keyring-backend file

# Create gentx
rewardchaind genesis gentx validator 50000000000stake \
  --chain-id "$CHAIN_ID" \
  --moniker $MONIKER \
  --commission-rate 0.1 \
  --commission-max-rate 0.2 \
  --commission-max-change-rate 0.01 \
  --min-self-delegation 1 \
  --keyring-backend file \
  --home "$REWARDCHAIN_HOME"

# Collect gentxs
rewardchaind genesis collect-gentxs --home "$REWARDCHAIN_HOME"

# Validate genesis
rewardchaind genesis validate-genesis --home "$REWARDCHAIN_HOME"
```

### 4. Configure Genesis (For Existing Chain)

If you're joining an existing network:

```bash
# Download the genesis file from the network
cd "$HOME/config"
wget https://your-network-genesis-url/genesis.json -O genesis.json

# Verify genesis file
rewardchaind genesis validate-genesis --home "$HOME"
```

## Configuration

### 1. Configure `config.toml`

Edit the P2P and RPC configuration:

```bash
nano "$HOME/config/config.toml"
```

**Key settings to configure:**

```bash
# Set moniker
sed -i 's/^moniker = .*/moniker = "validator-1"/' "$REWARDCHAIN_HOME/config/config.toml"

# Set external address (replace YOUR_SERVER_IP with actual IP)
sed -i 's/^external_address = .*/external_address = "157.175.215.244:26656"/' "$REWARDCHAIN_HOME/config/config.toml"

# Set seeds
# sed -i 's/^seeds = .*/seeds = "seed-node-1@ip:port,seed-node-2@ip:port"/' "$REWARDCHAIN_HOME/config/config.toml"

# Set persistent peers
# sed -i 's/^persistent_peers = .*/persistent_peers = "peer-1@ip:port,peer-2@ip:port"/' "$REWARDCHAIN_HOME/config/config.toml"

# Set RPC laddr to allow external access
sed -i 's/^laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' "$REWARDCHAIN_HOME/config/config.toml"

# Enable Prometheus
sed -i 's/^prometheus = .*/prometheus = true/' "$REWARDCHAIN_HOME/config/config.toml"

# Set Prometheus listen address
sed -i 's/^prometheus_listen_addr = .*/prometheus_listen_addr = ":26660"/' "$REWARDCHAIN_HOME/config/config.toml"
```

**Key settings to configure:**

```bash
# Set minimum gas prices
sed -i 's/^minimum-gas-prices = .*/minimum-gas-prices = "0.0001stake"/' "$REWARDCHAIN_HOME/config/app.toml"

# Enable API
sed -i '/\[api\]/,/\[/ s/^enable = .*/enable = true/' "$REWARDCHAIN_HOME/config/app.toml"

# Set API address
sed -i '/\[api\]/,/\[/ s/^address = .*/address = "tcp:\/\/0.0.0.0:1317"/' "$REWARDCHAIN_HOME/config/app.toml"

# Enable Swagger
sed -i '/\[api\]/,/\[/ s/^swagger = .*/swagger = true/' "$REWARDCHAIN_HOME/config/app.toml"

# Set gRPC address
sed -i '/\[grpc\]/,/\[/ s/^address = .*/address = "0.0.0.0:9090"/' "$REWARDCHAIN_HOME/config/app.toml"

# Set state-sync snapshot interval
sed -i 's/^snapshot-interval = .*/snapshot-interval = 1000/' "$REWARDCHAIN_HOME/config/app.toml"

# Set snapshot keep recent
sed -i 's/^snapshot-keep-recent = .*/snapshot-keep-recent = 2/' "$REWARDCHAIN_HOME/config/app.toml"

# Set pruning strategy
sed -i 's/^pruning = .*/pruning = "default"/' "$REWARDCHAIN_HOME/config/app.toml"

# Set pruning interval
sed -i 's/^pruning-interval = .*/pruning-interval = "10"/' "$REWARDCHAIN_HOME/config/app.toml"

# Set pruning keep recent
sed -i 's/^pruning-keep-recent = .*/pruning-keep-recent = "100"/' "$REWARDCHAIN_HOME/config/app.toml"

# Set pruning keep every
sed -i 's/^pruning-keep-every = .*/pruning-keep-every = "0"/' "$REWARDCHAIN_HOME/config/app.toml"

# Set pruning min retain blocks
sed -i 's/^min-retain-blocks = .*/min-retain-blocks = "0"/' "$REWARDCHAIN_HOME/config/app.toml"
```

### 3. Set Keyring Backend

For production, use file-based keyring (more secure than OS keyring on servers):

```bash
# The keyring backend is already set to "file" in the commands above
# Verify keyring location
ls -la "$REWARDCHAIN_HOME/keyring-file/"
```

## Systemd Service Setup

### 1. Create Systemd Service File

Create the systemd service file:

```bash
sudo nano /etc/systemd/system/rewardchaind.service
```

Add the following configuration:

```ini
[Unit]
Description=Reward Chain Daemon
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=ubuntu
Group=ubuntu
WorkingDirectory=/home/ubuntu
ExecStart=/usr/local/bin/rewardchaind start \
  --home /home/ubuntu/.rewardchain \
  --minimum-gas-prices 0.0001stake
Restart=always
RestartSec=3
LimitNOFILE=65535

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
# Remove or comment out ProtectHome=true if it causes issues
# ProtectHome=true
ReadWritePaths=/home/ubuntu/.rewardchain

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=rewardchaind

[Install]
WantedBy=multi-user.target
```

### 2. Reload and Enable Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable rewardchaind

# Start the service
sudo systemctl start rewardchaind

# Check status
sudo systemctl status rewardchaind

# View logs
sudo journalctl -u rewardchaind -f
```

### 3. Service Management Commands

```bash
# Start service
sudo systemctl start rewardchaind

# Stop service
sudo systemctl stop rewardchaind

# Restart service
sudo systemctl restart rewardchaind

# View logs
sudo journalctl -u rewardchaind -f

# View last 100 lines
sudo journalctl -u rewardchaind -n 100

# View logs since today
sudo journalctl -u rewardchaind --since today
```

## Firewall Configuration

Configure UFW (Uncomplicated Firewall) to allow necessary ports:

```bash
# Allow SSH (if not already configured)
sudo ufw allow 22/tcp

# Allow P2P port (required for node communication)
sudo ufw allow 26656/tcp

# Allow RPC port (optional, only if you need external RPC access)
sudo ufw allow 26657/tcp

# Allow gRPC port (optional, only if you need external gRPC access)
sudo ufw allow 9090/tcp

# Allow API port (optional, only if you need external API access)
sudo ufw allow 1317/tcp

# Allow Prometheus metrics port (optional)
sudo ufw allow 26660/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

**Security Note**: For production, consider restricting RPC, gRPC, and API access to specific IP addresses only:

```bash
# Example: Allow RPC only from specific IP
sudo ufw allow from YOUR_TRUSTED_IP to any port 26657
sudo ufw deny 26657/tcp
```

## Block Explorer Setup

This section covers setting up a [Ping.pub](https://ping.pub) block explorer for your Reward Chain node. The explorer provides a web interface to browse blocks, transactions, validators, and accounts.

### Prerequisites

Before setting up the explorer, ensure you have:

1. **Node.js** (v18 or later) - Install from [nodejs.org](https://nodejs.org/)
2. **Yarn** or **npm** - Package manager (Yarn is preferred)
3. **Git** - Should already be installed from prerequisites
4. **Running Reward Chain node** - Your node must be running with:
   - RPC endpoint accessible: `http://localhost:26657` (or your server IP)
   - REST API endpoint accessible: `http://localhost:1317` (or your server IP)

### Install Node.js and Yarn

If not already installed:

```bash
# Install Node.js (using NodeSource repository for latest LTS)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verify Node.js installation
node --version
npm --version

# Install Yarn (recommended)
npm install -g yarn

# Verify Yarn installation
yarn --version
```

### Configure CORS for Explorer Access

Before setting up the explorer, you need to enable CORS in your node configuration to allow the explorer to connect:

**1. Enable CORS in `config.toml`:**

```bash
nano /home/ubuntu/.rewardchain/config/config.toml
```

Add or update the `[rpc]` section:

```toml
[rpc]
cors_allowed_origins = ["*"]
```

**2. Enable CORS in `app.toml`:**

```bash
nano /home/ubuntu/.rewardchain/config/app.toml
```

Ensure the `[api]` section has:

```toml
[api]
enable = true
enabled-unsafe-cors = true
```

**3. Restart your node:**

```bash
sudo systemctl restart rewardchaind
```

### Automatic Setup (Recommended)

The easiest way to set up the explorer is using the provided setup script:

```bash
# Navigate to the explorer directory
# Replace ~/rewardchain with your actual repository path
cd ~/rewardchain/explorer

# Make the script executable
chmod +x setup.sh

# Run the setup script
./setup.sh setup
```

This script will:
1. Check for required dependencies (Node.js, Yarn/npm, Git)
2. Clone the Ping.pub explorer repository to `~/ping-pub-explorer`
3. Configure it for Reward Chain
4. Install dependencies
5. Start the development server

The explorer will be available at: **http://localhost:5173/rewardchain**

### Manual Setup

If you prefer to set up manually or the automatic script fails:

```bash
# 1. Clone the Ping.pub explorer
git clone https://github.com/ping-pub/explorer.git ~/ping-pub-explorer
cd ~/ping-pub-explorer

# 2. Install dependencies
yarn install
# Or if using npm (may require --legacy-peer-deps flag):
# npm install --legacy-peer-deps

# 3. Create chain configuration directory
mkdir -p chains/mainnet

# 4. Copy the Reward Chain configuration
# Replace ~/rewardchain with your actual repository path
cp ~/rewardchain/explorer/ping_pub_config.json chains/mainnet/rewardchain.json

# Rebuild with new config
yarn build    # or: npm run build

# Restart via PM2
pm2 start "serve -s dist -l 5173" --name reward-explorer
pm2 save
```

### Update Configuration for Remote Access

If you want to access the explorer from a remote machine, update the endpoints in the configuration file:

```bash
nano ~/ping-pub-explorer/chains/mainnet/rewardchain.json
```

Update the `api` and `rpc` addresses to use your server's IP:

```json
{
  "api": [
    {
      "address": "http://YOUR_SERVER_IP:1317",
      "provider": "local"
    }
  ],
  "rpc": [
    {
      "address": "http://YOUR_SERVER_IP:26657",
      "provider": "local"
    }
  ]
}
```

### Running Explorer as a Service (Optional)

For production, you may want to run the explorer as a systemd service:

```bash
sudo nano /etc/systemd/system/ping-pub-explorer.service
```

Add the following configuration:

```ini
[Unit]
Description=Ping.pub Block Explorer
After=network-online.target rewardchaind.service
Wants=network-online.target

[Service]
Type=simple
User=ubuntu
Group=ubuntu
WorkingDirectory=/home/ubuntu/ping-pub-explorer
ExecStart=/usr/bin/yarn dev
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=ping-pub-explorer

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ping-pub-explorer
sudo systemctl start ping-pub-explorer
sudo systemctl status ping-pub-explorer
```

### Firewall Configuration for Explorer

If you want to access the explorer from remote machines, allow the explorer port:

```bash
# Allow explorer port (5173)
sudo ufw allow 5173/tcp

# Check status
sudo ufw status
```

### Accessing the Explorer

Once running, open your browser and navigate to:

- **Local access**: http://localhost:5173/rewardchain
- **Remote access**: http://YOUR_SERVER_IP:5173/rewardchain

### Explorer Features

The Ping.pub explorer provides:

- **Blocks**: Browse recent blocks and block details
- **Transactions**: View transaction history and details
- **Validators**: See validator information and staking details
- **Accounts**: Check account balances and transaction history
- **Proposals**: View governance proposals (if applicable)
- **IBC**: Monitor IBC transfers and channels

### Troubleshooting Explorer Issues

**Transactions Not Showing:**

1. Verify API server is enabled: `enable = true` in `app.toml`
2. Verify CORS is enabled in both `config.toml` and `app.toml`
3. Check transaction indexing: `indexer = "kv"` in `config.toml`
4. **Restart your node** after making changes

**Port Already in Use:**

```bash
# Find what's using port 5173
sudo lsof -i :5173

# Kill the process or change the port
```

**Dependency Resolution Errors:**

If using npm and encountering dependency errors:

```bash
cd ~/ping-pub-explorer
npm install --legacy-peer-deps
```

**CORS Errors:**

1. Verify CORS is enabled in both configuration files
2. Restart your node after making changes
3. Check that the API server is enabled

**Node Not Responding:**

```bash
# Verify your node is running
curl http://localhost:26657/status
curl http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info
```

For more detailed troubleshooting, see the [Explorer Troubleshooting Guide](../explorer/TROUBLESHOOTING_TXS.md).

### Production Deployment

For production deployments:

1. **Build for production**:
   ```bash
   cd ~/ping-pub-explorer
   yarn build
   ```

2. **Use a production web server** (nginx, Apache) to serve the built files

3. **Update endpoints** in configuration to use production RPC/API URLs

4. **Configure HTTPS** for secure access

5. **Set up proper CORS** origins instead of `["*"]` for security

## Monitoring and Maintenance

### 1. Check Node Status

```bash
# Check if node is syncing
rewardchaind status --home /home/ubuntu/.rewardchain

# Check validator status
rewardchaind query staking validators --home /home/ubuntu/.rewardchain

# Check your validator
rewardchaind query staking validator $(rewardchaind keys show validator -a --bech val --home /home/ubuntu/.rewardchain --keyring-backend file) --home /home/ubuntu/.rewardchain
```

### 2. Monitor Logs

```bash
# Follow logs in real-time
sudo journalctl -u rewardchaind -f

# Check for errors
sudo journalctl -u rewardchaind | grep -i error

# Check sync status
sudo journalctl -u rewardchaind | grep -i "synced\|catching\|height"
```

### 3. Check Disk Usage

```bash
# Check disk usage
df -h

# Check data directory size
du -sh /home/ubuntu/.rewardchain/data

# Monitor disk usage over time
watch -n 60 'df -h /home/ubuntu/.rewardchain'
```

### 4. Set Up Log Rotation

Create a logrotate configuration:

```bash
sudo nano /etc/logrotate.d/rewardchaind
```

Add:

```
/var/log/rewardchaind/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 ubuntu ubuntu
}
```

### 5. Backup Strategy

**Important**: Regularly backup your validator keys and configuration.

```bash
# Create backup directory
mkdir -p ~/backups/rewardchain

# Backup keys (CRITICAL - store securely off-server)
tar -czf ~/backups/rewardchain/keys-$(date +%Y%m%d).tar.gz \
  /home/ubuntu/.rewardchain/keyring-file/

# Backup configuration
tar -czf ~/backups/rewardchain/config-$(date +%Y%m%d).tar.gz \
  /home/ubuntu/.rewardchain/config/

# Schedule automatic backups (add to crontab)
crontab -e
# Add: 0 2 * * * /path/to/backup-script.sh
```

## Troubleshooting

### Node Not Starting

```bash
# Check service status
sudo systemctl status rewardchaind

# Check logs for errors
sudo journalctl -u rewardchaind -n 50

# Common issues:
# 1. Port already in use
sudo netstat -tulpn | grep 26656

# 2. Insufficient permissions
ls -la /home/ubuntu/.rewardchain

# 3. Disk space full
df -h
```

### Node Not Syncing

```bash
# Check current block height
rewardchaind status --home /home/ubuntu/.rewardchain | jq .SyncInfo

# Check if connected to peers
rewardchaind status --home /home/ubuntu/.rewardchain | jq .NodeInfo

# Verify seed/peer configuration
cat /home/ubuntu/.rewardchain/config/config.toml | grep -A 5 "seeds\|persistent_peers"

# Check firewall
sudo ufw status
```

### Connection Issues

```bash
# Test P2P port
telnet YOUR_SERVER_IP 26656

# Check if port is listening
sudo netstat -tulpn | grep rewardchaind

# Verify external address in config.toml
cat /home/ubuntu/.rewardchain/config/config.toml | grep external_address
```

### High Memory Usage

```bash
# Check memory usage
free -h

# Monitor process memory
ps aux | grep rewardchaind

# Consider enabling pruning in app.toml
# See Configuration section above
```

### Validator Jailed

```bash
# Check validator status
rewardchaind query staking validator $(rewardchaind keys show validator -a --bech val --home /home/ubuntu/.rewardchain --keyring-backend file) --home /home/ubuntu/.rewardchain

# If jailed, check why
rewardchaind query slashing signing-info $(rewardchaind tendermint show-validator --home /home/ubuntu/.rewardchain) --home /home/ubuntu/.rewardchain

# Unjail (if appropriate)
rewardchaind tx slashing unjail \
  --from validator \
  --home /home/ubuntu/.rewardchain \
  --keyring-backend file \
  --chain-id "$CHAIN_ID"
```

### Service Keeps Restarting

```bash
# Check systemd logs
sudo journalctl -u rewardchaind -n 100

# Check for crash loops
sudo systemctl status rewardchaind

# Verify binary exists and is executable
ls -la /usr/local/bin/rewardchaind

# Check file permissions
ls -la /home/ubuntu/.rewardchain/
```

## Security Best Practices

1. **Firewall**: Only open necessary ports and restrict access where possible
2. **User Permissions**: Run the service as a non-root user
3. **Key Management**: Store validator keys securely and backup off-server
4. **Regular Updates**: Keep the system and binary updated
5. **Monitoring**: Set up monitoring and alerts for node health
6. **Backups**: Regularly backup keys and configuration
7. **SSH Security**: Use SSH keys instead of passwords, disable root login

## Next Steps

After your node is running:

1. **Monitor**: Set up monitoring dashboards (Prometheus/Grafana)
2. **Join Network**: If joining existing network, ensure you're connected to peers
3. **Bond Tokens**: Bond tokens to become an active validator
4. **Set Up Alerts**: Configure alerts for downtime, sync issues, etc.
5. **Documentation**: Document your specific configuration for your team

## Additional Resources

- [Build & Deploy Guide](./build-and-deploy.md) - General build and deployment information
- [Cosmos SDK Documentation](https://docs.cosmos.network)
- [Tendermint Documentation](https://docs.tendermint.com)

## Support

For issues specific to this chain, please refer to the main repository or contact the development team.

