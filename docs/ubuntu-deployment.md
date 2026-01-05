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
8. [Monitoring and Maintenance](#monitoring-and-maintenance)
9. [Troubleshooting](#troubleshooting)

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
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version
```

## Server Setup

### 1. Create a Dedicated User

For security best practices, create a dedicated user for running the node:

```bash
# Create user
sudo useradd -m -s /bin/bash rewardchain

# Add user to sudo group (optional, for administrative tasks)
sudo usermod -aG sudo rewardchain

# Switch to the new user
sudo su - rewardchain
```

### 2. Set Up Directory Structure

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
git clone https://github.com/your-username/reward-chain.git
cd reward-chain

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
cd /path/to/reward-chain
make install

# Build a release binary
mkdir -p build
go build -o ./build/rewardchaind ./cmd/reward-chaind

# Verify
./build/rewardchaind version
```

Transfer to the server:

```bash
# From your local machine
scp ./build/rewardchaind user@your-server-ip:/tmp/rewardchaind

# On the server
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
export MONIKER="your-validator-name"  # Your validator name
export HOME="$HOME/.rewardchain"

# Initialize the node
rewardchaind init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME"
```

### 2. Create or Import Validator Key

**Option A: Create a new key**

```bash
rewardchaind keys add validator \
  --home "$HOME" \
  --keyring-backend file
```

**Important**: Save the mnemonic phrase in a secure location. You'll need it to recover your key.

**Option B: Import existing key from mnemonic**

```bash
rewardchaind keys add validator \
  --recover \
  --home "$HOME" \
  --keyring-backend file
```

**Get your validator address:**

```bash
VAL_ADDR=$(rewardchaind keys show validator -a --home "$HOME" --keyring-backend file)
echo "Validator address: $VAL_ADDR"
```

### 3. Configure Genesis (For New Chain)

If you're starting a new chain:

```bash
# Add genesis account
rewardchaind genesis add-genesis-account "$VAL_ADDR" 1000000000stake --home "$HOME"

# Create validator gentx
rewardchaind genesis gentx validator 1000000stake \
  --chain-id "$CHAIN_ID" \
  --home "$HOME" \
  --keyring-backend file

# Collect gentxs
rewardchaind genesis collect-gentxs --home "$HOME"

# Validate genesis
rewardchaind genesis validate-genesis --home "$HOME"
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

```toml
# P2P Configuration
moniker = "your-validator-name"
external_address = "YOUR_SERVER_IP:26656"  # Your server's public IP

# Seed nodes (if joining existing network)
seeds = "seed-node-1@ip:port,seed-node-2@ip:port"

# Persistent peers
persistent_peers = "peer-1@ip:port,peer-2@ip:port"

# RPC Configuration
laddr = "tcp://0.0.0.0:26657"  # Allow external RPC access (or use 127.0.0.1 for local only)

# Enable Prometheus metrics
prometheus = true
prometheus_listen_addr = ":26660"
```

### 2. Configure `app.toml`

Edit the application configuration:

```bash
nano "$HOME/config/app.toml"
```

**Key settings to configure:**

```toml
# Minimum gas prices (required)
minimum-gas-prices = "0.0001stake"  # Adjust based on your token economics

# API Configuration
api {
  enable = true
  address = "tcp://0.0.0.0:1317"  # Or 127.0.0.1:1317 for local only
  swagger = true
}

# gRPC Configuration
grpc {
  address = "0.0.0.0:9090"  # Or 127.0.0.1:9090 for local only
}

# State Sync (optional, for faster syncing)
[state-sync]
snapshot-interval = 1000
snapshot-keep-recent = 2

# Pruning (for disk space management)
pruning = "default"  # Options: "default", "nothing", "everything", "custom"
pruning-interval = "10"
pruning-keep-recent = "100"
pruning-keep-every = "0"
pruning-min-retain-blocks = "0"
```

### 3. Set Keyring Backend

For production, use file-based keyring (more secure than OS keyring on servers):

```bash
# The keyring backend is already set to "file" in the commands above
# Verify keyring location
ls -la "$HOME/.rewardchain/keyring-file/"
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
User=rewardchain
Group=rewardchain
WorkingDirectory=/home/rewardchain
ExecStart=/usr/local/bin/rewardchaind start \
  --home /home/rewardchain/.rewardchain \
  --minimum-gas-prices 0.0001stake
Restart=always
RestartSec=3
LimitNOFILE=65535

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ProtectHome=true
ReadWritePaths=/home/rewardchain/.rewardchain

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

## Monitoring and Maintenance

### 1. Check Node Status

```bash
# Check if node is syncing
rewardchaind status --home /home/rewardchain/.rewardchain

# Check validator status
rewardchaind query staking validators --home /home/rewardchain/.rewardchain

# Check your validator
rewardchaind query staking validator $(rewardchaind keys show validator -a --bech val --home /home/rewardchain/.rewardchain --keyring-backend file) --home /home/rewardchain/.rewardchain
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
du -sh /home/rewardchain/.rewardchain/data

# Monitor disk usage over time
watch -n 60 'df -h /home/rewardchain/.rewardchain'
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
    create 0640 rewardchain rewardchain
}
```

### 5. Backup Strategy

**Important**: Regularly backup your validator keys and configuration.

```bash
# Create backup directory
mkdir -p ~/backups/rewardchain

# Backup keys (CRITICAL - store securely off-server)
tar -czf ~/backups/rewardchain/keys-$(date +%Y%m%d).tar.gz \
  /home/rewardchain/.rewardchain/keyring-file/

# Backup configuration
tar -czf ~/backups/rewardchain/config-$(date +%Y%m%d).tar.gz \
  /home/rewardchain/.rewardchain/config/

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
ls -la /home/rewardchain/.rewardchain

# 3. Disk space full
df -h
```

### Node Not Syncing

```bash
# Check current block height
rewardchaind status --home /home/rewardchain/.rewardchain | jq .SyncInfo

# Check if connected to peers
rewardchaind status --home /home/rewardchain/.rewardchain | jq .NodeInfo

# Verify seed/peer configuration
cat /home/rewardchain/.rewardchain/config/config.toml | grep -A 5 "seeds\|persistent_peers"

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
cat /home/rewardchain/.rewardchain/config/config.toml | grep external_address
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
rewardchaind query staking validator $(rewardchaind keys show validator -a --bech val --home /home/rewardchain/.rewardchain --keyring-backend file) --home /home/rewardchain/.rewardchain

# If jailed, check why
rewardchaind query slashing signing-info $(rewardchaind tendermint show-validator --home /home/rewardchain/.rewardchain) --home /home/rewardchain/.rewardchain

# Unjail (if appropriate)
rewardchaind tx slashing unjail \
  --from validator \
  --home /home/rewardchain/.rewardchain \
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
ls -la /home/rewardchain/.rewardchain/
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

