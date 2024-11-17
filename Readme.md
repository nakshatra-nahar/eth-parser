# Ethereum Transaction Parser

## Overview

The Ethereum Transaction Parser is a Go application designed to parse the Ethereum blockchain and monitor transactions involving subscribed addresses. It allows users to:

- **Subscribe to Ethereum addresses** they want to monitor.
- **Retrieve inbound and outbound transactions** for those addresses.
- **Check the last parsed block number** to monitor parsing progress.

The application interacts with the Ethereum blockchain via JSON-RPC calls to a public node and uses in-memory storage for simplicity. It's built with extensibility in mind, allowing for future integration with persistent storage solutions or notification services.

## Features

- **Address Subscription:** Subscribe to Ethereum addresses to monitor transactions.
- **Transaction Retrieval:** Fetch inbound and outbound transactions for subscribed addresses.
- **Current Block Tracking:** Retrieve the last parsed block number.
- **HTTP API Endpoints:** Interact with the parser using a simple HTTP API.
- **In-Memory Data Storage:** Utilizes in-memory storage for quick access (can be extended to persistent storage).

## Requirements

- **Go 1.16 or higher:** Install Go from the [official website](https://golang.org/dl/).
- **Internet Connection:** Required for interacting with the Ethereum blockchain.
- **Access to Public Ethereum Node:** The application uses `https://ethereum-rpc.publicnode.com` for JSON-RPC calls.

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/nakshatra-nahar/eth-parser.git
cd eth-parser
```

### 2. Verify Go Installation

Ensure Go is installed and accessible:

```bash
go version
```

### 3. Review Project Structure

The project consists of the following files:

- `main.go`: Contains the `main` function and sets up the HTTP server.
- `parser.go`: Implements the `Parser` interface and the `EthereumParser` struct.
- `transaction.go`: Defines the `Transaction` struct.
- `rpc.go`: Contains functions for interacting with the Ethereum JSON-RPC API.
- `utils.go`: Provides utility functions like `hexToDecimal` and `calculateTransactionFee`.

## Usage

### Running the Application

You can run the application directly using `go run` or build it into an executable.

#### Option 1: Run Directly

```bash
go run *.go
```

#### Option 2: Build and Execute

Build the application:

```bash
go build -o eth-parser
```

Run the executable:

```bash
./eth-parser
```

Upon running, you should see output similar to:

```
Starting parser from block: [block number]
Server is listening on port 8080
```

### API Endpoints

The application exposes several HTTP endpoints:

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

### Testing the Application

1. **Subscribe to an Address**

   Replace `0xYourEthereumAddress` with the Ethereum address you wish to monitor.

   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"address":"0xYourEthereumAddress"}' http://localhost:8080/subscribe
   ```

2. **Check Current Block**

   Verify that the parser is processing blocks:

   ```bash
   curl http://localhost:8080/currentBlock
   ```

3. **Fetch Transactions**

   After a few minutes (to allow the parser to process blocks), retrieve transactions for the subscribed address:

   ```bash
   curl 'http://localhost:8080/transactions?address=0xYourEthereumAddress'
   ```

4. **Get Transaction Details**

   Get details of a specific transaction by its hash:

   ```bash
   curl 'http://localhost:8080/transaction?hash=0xYourTransactionHash'
   ```

### Configuration

- **Starting Block:** By default, the parser starts from 10,000 blocks before the latest block. You can adjust this by modifying the `N` constant in `parser.go`.

  ```go
  const N int64 = 10000 // Number of blocks before the latest block
  ```

- **Delay Between Block Processing:** To prevent rate limiting, the parser includes a delay between processing blocks. Adjust the delay in the `StartParsing` method in `parser.go`.

  ```go
  time.Sleep(200 * time.Millisecond) // Adjust the delay as needed
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

- **main.go:** Initializes the parser and sets up the HTTP server.
- **parser.go:** Contains the parser logic and interfaces.
- **transaction.go:** Defines the data structure for transactions.
- **rpc.go:** Handles JSON-RPC communication with the Ethereum node.
- **utils.go:** Provides utility functions for data conversion and calculation.

## Troubleshooting

- **Empty Transaction Array (`[]`):**

  - **Cause:** The parser may not have processed blocks containing transactions for your address yet.
  - **Solution:** Wait for the parser to process more blocks or increase the `N` value to start from an earlier block.

- **Connection Errors:**

  - **Cause:** Network issues or problems connecting to the Ethereum node.
  - **Solution:** Check your internet connection and ensure the Ethereum JSON-RPC endpoint is accessible.

- **Rate Limiting:**

  - **Cause:** Exceeding request limits on the public Ethereum node.
  - **Solution:** Increase the delay between requests in the `StartParsing` method.

- **Incorrect Address Format:**

  - Ensure that Ethereum addresses are in the correct format (checksummed or all lowercase). The parser normalizes addresses to lowercase.

## Extending the Application

- **Persistent Storage:** Implement a database to store subscribed addresses and transactions persistently.
- **Notification Services:** Integrate with messaging services (e.g., email, SMS, push notifications) to alert users of new transactions.
- **Web Interface:** Develop a frontend to interact with the parser visually.
- **Additional Blockchain Support:** Extend support to other blockchains with similar JSON-RPC interfaces.

---