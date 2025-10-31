# Immutable Blockchain Log

This project implements an immutable log system using Hyperledger Fabric for on-chain storage and PostgreSQL for off-chain storage. Log entries are hashed and stored on the blockchain for immutability, with content stored off-chain.

## Prerequisites

- Go 1.23
- Docker and Docker Compose
- Hyperledger Fabric binaries (will be installed automatically by ./network-up.sh)

## Setup

1. Clone the repository and navigate to the project root directory.

2. Run the blockchain network setup:
   ```sh
   ./network-up.sh
   ```
   This script will install binaries if needed, start the Fabric network, create a channel, and deploy the chaincode. Wait for it to complete (up to few minutes).

3. Start the log dashboard and API gateway:
   ```sh
   ./logger-up.sh
   ```
   This script will:
   - Build the log-client API gateway (backend)
   - Install dependencies and build the log-dashboard (frontend)
   - Start both the API gateway (port 3001) and the web dashboard (port 3000)
   
   Access the dashboard at `http://localhost:3000`

## Usage

### Using the Web Dashboard

The web dashboard provides a user-friendly interface to interact with the blockchain log system:

1. **Settings**: Configure the log file path to monitor
2. **Logs**: View all log entries with validation status

The dashboard automatically connects to the API gateway running on port 3001.

### Using Command Line (Deprecated - use web dashboard)
#### Writing Logs

1. Change to the log-client directory:
   ```sh
   cd log-client
   ```

2. Tidy Go modules:
   ```sh
   go mod tidy
   ```

3. Run the write-log command to monitor a text file for new lines:
   ```sh
   go run cmd/write-log/main.go <filename> <client-name>
   ```
   Example:
   ```sh
   go run cmd/write-log/main.go test.txt client1
   ```
   This starts monitoring [`log-client/test.txt`](log-client/test.txt) for new lines. Each new line is written to the off-chain database and a corresponding asset (with hash) is created on the blockchain.

#### Reading Logs

1. From the log-client directory, run the read-log command:
   ```sh
   go run cmd/read-log/main.go [client-name-filter]
   ```
   Example:
   ```sh
   go run cmd/read-log/main.go client1
   ```
   This retrieves logs from the blockchain and off-chain storage, validates them using hashes, and displays the results. If no filter is provided, all logs are shown.

## Code Structure
The project is organized into the following key directories and files:

- **Root Directory**: Contains setup scripts and configuration.
  - [`install_binary.sh`](install_binary.sh ): Installs Hyperledger Fabric binaries if not present.
  - [`network-up.sh`](network-up.sh ): Sets up the Fabric network, creates channels, and deploys chaincode.
  - [`compose-test-net.yaml`](compose-test-net.yaml ): Docker Compose configuration for the network.
  - [`config`](config ): Fabric configuration files (e.g., [`configtx.yaml`](config/configtx.yaml ), [`core.yaml`](config/core.yaml )).
  - [`channel-artifacts`](channel-artifacts ): Generated channel artifacts (e.g., [`mychannel.block`](channel-artifacts/mychannel.block )).
  - [`organizations`](organizations ): Cryptographic materials for peers and orderers.

- **bin/**: Directory for installed Hyperledger Fabric binaries (e.g., `peer`, `orderer`, `cryptogen`).

- **chaincode-go/**: Go-based chaincode implementation.
  - [`assetTransfer.go`](chaincode-go/assetTransfer.go ): Main entry point for the chaincode.
  - [`chaincode/smartcontract.go`](chaincode-go/chaincode/smartcontract.go ): Defines the [`SmartContract`](chaincode-go/chaincode/smartcontract.go ) struct with methods like [`CreateAsset`](chaincode-go/chaincode/smartcontract.go ), [`AssetExists`](chaincode-go/chaincode/smartcontract.go ), and [`GetAllAssets`](chaincode-go/chaincode/smartcontract.go ).

- **log-client/**: Go client application for interacting with the blockchain and off-chain storage.
  - `cmd/`: Command-line interfaces.
    - [`gateway/main.go`](log-client/cmd/gateway/main.go ): REST API gateway service using Gin framework. Provides HTTP endpoints for log management, file monitoring configuration, and log retrieval with validation.
    - [`write-log/main.go`](log-client/cmd/write-log/main.go ): Monitors a file for new lines, writes to PostgreSQL, and creates blockchain assets.
    - [`read-log/main.go`](log-client/cmd/read-log/main.go ): Retrieves and validates logs from blockchain and database.
  - `internal/`: Internal packages.
    - [`grpc-connection.go`](log-client/internal/grpc-connection.go ): Manages gRPC connections to Fabric Gateway.
    - [`database.go`](log-client/internal/database.go ): Initializes PostgreSQL connection using GORM.
    - [`log-entry.go`](log-client/internal/log-entry.go ): Defines [`LogEntry`](log-client/internal/log-entry.go ) struct with methods like [`Hash`](log-client/internal/log-entry.go ), [`ValidateHash`](log-client/internal/log-entry.go ), [`LoadFromDB`](log-client/internal/log-entry.go ), and [`WriteToDB`](log-client/internal/log-entry.go ).
    - [`utils.go`](log-client/internal/utils.go ): File watching utility with [`WatchFile`](log-client/internal/utils.go ).
    - [`constants.go`](log-client/internal/constants.go ): Constants for MSP ID, crypto paths, endpoints, etc.

- **log-dashboard/**: React-based web dashboard for the log system.
  - Built with Vite, TypeScript, and TanStack Router
  - `src/`: Source code for the dashboard.
    - `components/`: React components including dashboard, logs viewer, settings, and UI components.
    - `routes/`: Page routes for navigation (dashboard, logs, settings).
  - Provides a web interface to configure log monitoring, view logs, and monitor system activity.
