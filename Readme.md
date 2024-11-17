# Ethereum Transaction Parser

## Overview

The Ethereum Transaction Parser is a Go application that parses the Ethereum blockchain to monitor transactions for subscribed addresses. It provides features to:

- **Subscribe to Ethereum addresses** and monitor their transactions.
- **Retrieve inbound and outbound transactions** for subscribed addresses.
- **Check the last parsed block number** to monitor the parsing progress.

The parser interacts with the Ethereum blockchain using JSON-RPC calls and employs in-memory storage for simplicity, making it easy to extend with persistent storage or notification services in the future.

## Features

- **Address Subscription:** Monitor transactions for specified Ethereum addresses.
- **Transaction Retrieval:** Fetch inbound/outbound transactions for subscribed addresses.
- **Current Block Tracking:** Check the latest block processed by the parser.
- **HTTP API Endpoints:** Simple HTTP API for interaction.
- **In-Memory Data Storage:** Quick access to data, with the option to extend to persistent storage.

## Requirements

- **Go 1.16 or higher:** [Download Go](https://golang.org/dl/)
- **Internet Connection:** Needed for connecting to the Ethereum network.
- **Ethereum Node Access:** Uses `https://ethereum-rpc.publicnode.com` for JSON-RPC interactions.

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/nakshatra-nahar/eth-parser.git
cd eth-parser
```

### 2. Verify Go Installation

Make sure Go is installed:

```bash
go version
```

### 3. Review the Project Structure

- `main.go`: Initializes the application and sets up the HTTP server.
- `parser.go`: Contains the Ethereum parser logic.
- `transaction.go`: Defines the `Transaction` struct.
- `rpc.go`: Handles communication with the Ethereum JSON-RPC API.
- `utils.go`: Provides utility functions for data conversion and fee calculation.

## Usage

### Running the Application

You can run the parser using `go run` or build it into an executable.

#### Option 1: Run Directly

```bash
go run *.go
```

#### Option 2: Build and Run

Build the executable:

```bash
go build -o eth-parser
```

Run the executable:

```bash
./eth-parser
```

Expected output:

```
Starting parser from block: [block number]
Server is listening on port 8080
```

### API Endpoints

#### 1. Subscribe to an Address

**Endpoint:**

```
POST /subscribe
```

**Request Body:**

```json
{
  "address": "0xYourEthereumAddress"
}
```

**Example:**

```bash
curl -X POST -H "Content-Type: application/json" -d '{"address":"0xYourEthereumAddress"}' http://localhost:8080/subscribe
```

#### 2. Get Transactions for an Address

**Endpoint:**

```
GET /transactions?address=0xYourEthereumAddress
```

**Example:**

```bash
curl 'http://localhost:8080/transactions?address=0xYourEthereumAddress'
```

#### 3. Get Current Block Number

**Endpoint:**

```
GET /currentBlock
```

**Example:**

```bash
curl http://localhost:8080/currentBlock
```

#### 4. Get Transaction Details by Hash

**Endpoint:**

```
GET /transaction?hash=0xYourTransactionHash
```

**Example:**

```bash
curl 'http://localhost:8080/transaction?hash=0xYourTransactionHash'
```

### Configuration

- **Starting Block:** The parser starts from 10,000 blocks before the latest block by default. Modify the `N` constant in `parser.go` to adjust this.

  ```go
  const N int64 = 10000 // Number of blocks before the latest block
  ```

- **Delay Between Block Processing:** To avoid rate limiting, the parser has a delay between processing blocks. Adjust it in the `StartParsing` method in `parser.go`.

  ```go
  time.Sleep(200 * time.Millisecond) // Modify the delay as needed
  ```

## Project Structure

```
eth-parser/
├── main.go
├── parser.go
├── transaction.go
├── rpc.go
└── utils.go
```

- **main.go:** Initializes the parser and starts the HTTP server.
- **parser.go:** Implements the parsing logic and methods for address subscription.
- **transaction.go:** Defines the structure of a transaction.
- **rpc.go:** Handles JSON-RPC API calls to the Ethereum node.
- **utils.go:** Contains utility functions for hex-to-decimal conversion and fee calculation.

## Troubleshooting

- **Empty Transaction Array (`[]`):**
  - **Cause:** The parser may not have processed blocks containing relevant transactions yet.
  - **Solution:** Wait for more blocks to be processed or adjust the `N` value to start from an earlier block.

- **Connection Issues:**
  - **Cause:** Network errors or issues with the Ethereum node.
  - **Solution:** Check your network and ensure the Ethereum JSON-RPC endpoint is reachable.

- **Rate Limiting:**
  - **Cause:** Making too many requests to the Ethereum node too quickly.
  - **Solution:** Increase the delay between requests in the `StartParsing` method.

- **Incorrect Address Format:**
  - Ethereum addresses should be checksummed or all lowercase. The parser normalizes addresses to lowercase.

## Extending the Application

- **Persistent Storage:** Replace the in-memory storage with a database (e.g., PostgreSQL, Redis) for long-term data persistence.
- **Notification System:** Integrate with services like Twilio, Slack, or AWS SNS for transaction alerts.
- **Web Interface:** Build a frontend to display monitored transactions and parsing progress.
- **Multi-Blockchain Support:** Add support for other blockchains with similar JSON-RPC APIs.
