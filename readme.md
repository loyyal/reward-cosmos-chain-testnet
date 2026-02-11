# rewardchain
**rewardchain** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Prerequisites

Before you begin, ensure you have the following dependencies installed on your system:

### Required Dependencies

#### 1. Go (v1.24 or later)

This project requires **Go 1.24.0 or later** (as specified in `go.mod`). Any patch version of Go 1.24.x will work.

**Installation:**

**macOS (using Homebrew):**
```bash
brew install go@1.24
# Or install the latest Go version
brew install go
```

**Linux:**
```bash
# Download and install the latest Go 1.24.x version
# Check https://go.dev/dl/ for the latest version
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
```

**Or using a package manager (Ubuntu/Debian):**
```bash
sudo apt-get update
sudo apt-get install golang-go
```

**Windows:**
Download the latest Go 1.24.x installer from [https://go.dev/dl/](https://go.dev/dl/)

**Verify installation:**
```bash
go version
# Should output: go version go1.24.x ... (where x is any patch version)
```

**Set up Go environment variables** (if not already set):
```bash
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Add these to your `~/.bashrc` or `~/.zshrc` to make them permanent.

#### 2. Ignite CLI

**Installation:**

**macOS/Linux:**
```bash
curl https://get.ignite.com/cli! | bash
```

**Or using Homebrew (macOS):**
```bash
brew install ignite
```

**Or using npm:**
```bash
npm install -g @ignite/cli
```

**Verify installation:**
```bash
ignite version
```

#### 3. Git

**Installation:**

**macOS:**
```bash
brew install git
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get update
sudo apt-get install git
```

**Linux (Fedora/RHEL):**
```bash
sudo dnf install git
```

**Windows:**
Download from [https://git-scm.com/download/win](https://git-scm.com/download/win)

**Verify installation:**
```bash
git --version
```

#### 4. GNU Make (optional but recommended)

**Installation:**

**macOS:**
```bash
brew install make
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install build-essential
```

**Linux (Fedora/RHEL):**
```bash
sudo dnf groupinstall "Development Tools"
```

**Verify installation:**
```bash
make --version
```

### Optional Dependencies

#### 5. buf (for Protocol Buffer generation)

The project uses `buf` v1.34.0 for Protocol Buffer code generation.

**Installation:**

**macOS/Linux:**
```bash
brew install bufbuild/buf/buf
```

**Or using the official installer:**
```bash
# macOS
brew tap bufbuild/buf
brew install buf

# Linux
# Download from https://github.com/bufbuild/buf/releases
```

**Verify installation:**
```bash
buf --version
```

#### 6. golangci-lint (for code linting)

The project uses **golangci-lint v1.61.0** (as specified in `Makefile`).

**Installation:**
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
```

**Verify installation:**
```bash
golangci-lint --version
```

#### 7. jq (optional, helpful for JSON processing)

**Installation:**

**macOS:**
```bash
brew install jq
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install jq
```

**Linux (Fedora/RHEL):**
```bash
sudo dnf install jq
```

**Verify installation:**
```bash
jq --version
```

### Version Summary

| Dependency | Required Version | Purpose |
|------------|------------------|---------|
| Go | 1.24.0 or later (1.24.x) | Core programming language |
| Ignite CLI | Latest stable | Blockchain scaffolding and development |
| Git | Latest | Version control |
| GNU Make | Latest | Build automation |
| buf | v1.34.0 | Protocol Buffer tooling |
| golangci-lint | v1.61.0 | Code linting |

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

## Build & deploy (binary / node ops)

See `docs/build-and-deploy.md`.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Additionally, Ignite CLI offers both Vue and React options for frontend scaffolding:

For a Vue frontend, use: `ignite scaffold vue`
For a React frontend, use: `ignite scaffold react`
These commands can be run within your scaffolded blockchain project. 


For more information see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/username/rewardchain@latest! | sudo bash
```
`username/rewardchain` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
