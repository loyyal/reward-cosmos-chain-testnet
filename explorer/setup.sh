#!/bin/bash

# ============================================================================
# Ping.pub Explorer Setup Script
# Sets up the Ping.pub block explorer for the testnet
# ============================================================================

set -e

# Get script directory (works even when called from different directory)
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Configuration
EXPLORER_DIR="$HOME/ping-pub-explorer"
CHAIN_CONFIG_FILE="$SCRIPT_DIR/ping_pub_config.json"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check dependencies
check_dependencies() {
    log_info "Checking dependencies..."
    
    if ! command -v node &> /dev/null; then
        log_error "Node.js is required. Install from https://nodejs.org/"
        exit 1
    fi
    
    if ! command -v yarn &> /dev/null; then
        log_warning "Yarn not found, using npm..."
        USE_NPM=true
    else
        USE_NPM=false
    fi
    
    if ! command -v git &> /dev/null; then
        log_error "Git is required."
        exit 1
    fi
    
    log_success "Dependencies checked"
}

# Clone Ping.pub repository
clone_explorer() {
    log_info "Setting up Ping.pub explorer..."
    
    if [ -d "$EXPLORER_DIR" ]; then
        log_warning "Explorer directory exists. Updating..."
        cd $EXPLORER_DIR
        git pull origin main || git pull origin master
    else
        log_info "Cloning Ping.pub explorer..."
        git clone https://github.com/ping-pub/explorer.git $EXPLORER_DIR
    fi
    
    log_success "Ping.pub repository ready"
}

# Configure chain
configure_chain() {
    log_info "Configuring chain in explorer..."
    
    cd $EXPLORER_DIR
    
    # Create chains directory if it doesn't exist
    mkdir -p chains/mainnet
    
    # Copy our chain configuration
    cp $CHAIN_CONFIG_FILE chains/mainnet/rewardchain.json
    
    # Create a simplified configuration for the explorer
    cat > chains/mainnet/rewardchain.json << 'EOF'
{
  "chain_name": "rewardchain",
  "api": ["http://157.175.215.244:1317"],
  "rpc": ["http://157.175.215.244:26657"],
  "sdk_version": "0.50.11",
  "coin_type": "118",
  "min_tx_fee": "5000",
  "addr_prefix": "reward",
  "logo": "/logos/cosmos.svg",
  "assets": [{
    "base": "stake",
    "symbol": "STAKE",
    "exponent": "6",
    "logo": "/logos/cosmos.svg"
  }]
}
EOF
    
    log_success "Chain configured"
}

# Install dependencies
install_dependencies() {
    log_info "Installing explorer dependencies..."
    
    cd $EXPLORER_DIR
    
    if [ "$USE_NPM" = true ]; then
        log_info "Using npm with --legacy-peer-deps to resolve dependency conflicts..."
        npm install --legacy-peer-deps
    else
        yarn install
    fi
    
    log_success "Dependencies installed"
}

# Start explorer
start_explorer() {
    log_info "Starting Ping.pub explorer..."
    
    cd $EXPLORER_DIR
    
    echo ""
    echo "============================================"
    echo "   Starting Explorer"
    echo "============================================"
    echo ""
    echo "The explorer will be available at:"
    echo "  http://localhost:5173"
    echo ""
    echo "Make sure your testnet is running on:"
    echo "  RPC: http://localhost:26657"
    echo "  API: http://localhost:1317"
    echo ""
    echo "Press Ctrl+C to stop the explorer"
    echo ""
    
    if [ "$USE_NPM" = true ]; then
        npm run dev
    else
        yarn dev
    fi
}

# Print manual setup instructions
print_manual_setup() {
    echo ""
    echo "============================================"
    echo "   Manual Ping.pub Setup Instructions"
    echo "============================================"
    echo ""
    echo "If the automatic setup fails, follow these steps:"
    echo ""
    echo "1. Clone the Ping.pub explorer:"
    echo "   git clone https://github.com/ping-pub/explorer.git"
    echo "   cd explorer"
    echo ""
    echo "2. Install dependencies:"
    echo "   yarn install  # or npm install --legacy-peer-deps"
    echo "   (Note: npm may require --legacy-peer-deps flag to resolve dependency conflicts)"
    echo ""
    echo "3. Create chain configuration:"
    echo "   mkdir -p chains/mainnet"
    echo "   cp $CHAIN_CONFIG_FILE chains/mainnet/rewardchain.json"
    echo ""
    echo "4. Start the development server:"
    echo "   yarn dev  # or npm run dev"
    echo ""
    echo "5. Access the explorer at:"
    echo "   http://localhost:5173/rewardchain"
    echo ""
    echo "============================================"
    echo "   Alternative: Use Hosted Ping.pub"
    echo "============================================"
    echo ""
    echo "You can also use the hosted version:"
    echo "  https://ping.pub"
    echo ""
    echo "To add your local testnet:"
    echo "1. Go to https://ping.pub"
    echo "2. Click 'Add Chain' in the sidebar"
    echo "3. Enter your local RPC: http://localhost:26657"
    echo "4. The explorer will auto-detect chain info"
    echo ""
    echo "Note: For local chains, you may need to run"
    echo "a CORS proxy or configure your node for CORS."
    echo ""
}

# Enable CORS on the node
enable_cors_instructions() {
    echo ""
    echo "============================================"
    echo "   Enable CORS for Explorer Access"
    echo "============================================"
    echo ""
    echo "To allow the explorer to connect to your node,"
    echo "update the following configuration files:"
    echo ""
    echo "1. Edit config.toml:"
    echo "   ~/.rewardchain/config/config.toml"
    echo ""
    echo '   [rpc]'
    echo '   cors_allowed_origins = ["*"]'
    echo ""
    echo "2. Edit app.toml:"
    echo "   ~/.rewardchain/config/app.toml"
    echo ""
    echo '   [api]'
    echo '   enable = true'
    echo '   enabled-unsafe-cors = true'
    echo ""
    echo "3. Restart your node"
    echo ""
}

# Main
main() {
    echo "============================================"
    echo "   Ping.pub Explorer Setup"
    echo "============================================"
    echo ""
    
    case "${1:-setup}" in
        setup)
            check_dependencies
            clone_explorer
            configure_chain
            install_dependencies
            enable_cors_instructions
            start_explorer
            ;;
        manual)
            print_manual_setup
            enable_cors_instructions
            ;;
        cors)
            enable_cors_instructions
            ;;
        *)
            echo "Usage: $0 [setup|manual|cors]"
            echo ""
            echo "  setup  - Automatically setup and start explorer"
            echo "  manual - Print manual setup instructions"
            echo "  cors   - Print CORS configuration instructions"
            exit 1
            ;;
    esac
}

main "$@"