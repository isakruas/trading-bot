# Trading Bot CLI

A command-line application for interacting with cryptocurrency exchanges through a hexagonal (ports & adapters) architecture.  
Currently supports Foxbit (Brazilian exchange); others (Binance, Coinbase) can be added via new adapters.

---

## Table of Contents

1. [Features](#features)  
2. [Architecture](#architecture)  
3. [Prerequisites](#prerequisites)  
4. [Installation](#installation)  
5. [Configuration](#configuration)  
6. [Usage](#usage)  
   - [Global Help](#global-help)  
   - [fetch-markets](#fetch-markets)  
   - [fetch-order-book](#fetch-order-book)  
   - [place-order](#place-order)  
   - [cancel-order](#cancel-order)  
   - [list-active-orders](#list-active-orders)  
   - [get-order](#get-order)  
7. [Error Handling](#error-handling)  
8. [Extending to Other Exchanges](#extending-to-other-exchanges)  
9. [License](#license)  

---

## Features

- List available trading markets  
- Retrieve top-of-book bids and asks  
- Place new limit orders (post-only, GTC)  
- Cancel existing orders  
- List active orders by market  
- Fetch details of a single order  
- Human-readable tabular display  
- Structured error formatting for Foxbit’s JSON-style errors  
- Clean separation of concerns (use cases, domain, adapters, CLI)

---

## Architecture

This CLI is implemented using a Hexagonal (Ports & Adapters) pattern:

- **Domain** (`internal/domain`)  
  - `model` — core data structures (`Market`, `Order`, `OrderBook`)  
  - `service` — port interface `Exchange` defining available operations  

- **Application** (`internal/application/usecase`)  
  - Encapsulates each use case (e.g. `FetchMarkets`, `PlaceOrder`)  
  - Depends only on the `Exchange` interface  

- **Infrastructure** (`internal/infrastructure`)  
  - `exchange/foxbit` — adapter implementing `Exchange` via Foxbit REST API  
  - `httputil` — HTTP client helper for signing, sending, and parsing requests  

- **Interfaces** (`internal/interfaces/cli`)  
  - CLI entry point & command dispatch (`runner.go`)  
  - Table display helpers (`display.go`)  
  - Error formatter (`error_formatter.go`)

---

## Prerequisites

- Go 1.18+  
- A Foxbit account with API Key & Secret  

---

## Installation

```bash
# Clone the repository
git clone https://github.com/isakruas/trading-bot.git
cd trading-bot

# Build the CLI binary
make build
```

---

## Configuration

Set the following environment variables before running the CLI:

```bash
export FOXBIT_API_KEY="your_foxbit_api_key"
export FOXBIT_API_SECRET="your_foxbit_api_secret"
```

---

## Usage

General syntax:

```bash
trading-bot <command> [options]
```

To display global help:

```bash
trading-bot help
```

### fetch-markets

List all available trading markets.

```
Usage: trading-bot fetch-markets [--exchange foxbit|binance|coinbase]
```

Example:

```bash
trading-bot fetch-markets --exchange foxbit
```

### fetch-order-book

Retrieve the top-of-book bids and asks for a specific market.

```
Usage: trading-bot fetch-order-book --market SYMBOL [--depth N] [--exchange foxbit]
```

Options:

- `--market` (required) — market symbol (e.g. `BTCBRL`)  
- `--depth` — number of bid/ask levels (default: `10`)  
- `--exchange` — adapter name (default: `foxbit`)

Example:

```bash
trading-bot fetch-order-book --market BTCBRL --depth 5
```

### place-order

Place a new limit order (post-only, GTC).

```
Usage: trading-bot place-order --market SYMBOL --quantity QTY --price PRICE [--side buy|sell] [--exchange foxbit]
```

Options:

- `--market` (required) — market symbol  
- `--quantity` (required) — order quantity  
- `--price` (required) — order price  
- `--side` — `buy` or `sell` (default: `buy`)  
- `--exchange` — adapter name (default: `foxbit`)

Example:

```bash
trading-bot place-order --market BTCBRL --quantity 0.005 --price 100000.00 --side sell
```

### cancel-order

Cancel an existing order by ID.

```
Usage: trading-bot cancel-order --order-id ID [--exchange foxbit]
```

Options:

- `--order-id` (required) — ID of the order to cancel  

Example:

```bash
trading-bot cancel-order --order-id a1b2c3d4
```

### list-active-orders

List all active (OPEN) orders for a given market.

```
Usage: trading-bot list-active-orders --market SYMBOL [--exchange foxbit]
```

Options:

- `--market` (required) — market symbol  

Example:

```bash
trading-bot list-active-orders --market BTCBRL
```

### get-order

Fetch details of a single order by ID.

```
Usage: trading-bot get-order --order-id ID [--exchange foxbit]
```

Options:

- `--order-id` (required) — order ID  

Example:

```bash
trading-bot get-order --order-id a1b2c3d4
```

---

## Error Handling

- If an HTTP or JSON error occurs, the CLI attempts to parse Foxbit’s standard error payload:
  ```json
  {
    "error": {
      "message": "Invalid symbol",
      "code": 400,
      "details": ["market_symbol: not found"]
    }
  }
  ```
- Errors are rendered in a neat tabular form (`TYPE | CODE | MESSAGE | DETAIL`) to `stderr`.  
- Unknown or non-JSON errors are printed as raw text.

---

## Extending to Other Exchanges

1. Implement a new adapter under `internal/infrastructure/exchange/` that satisfies the `service.Exchange` interface.  
2. Register it in `mustInitExchange` (in `internal/interfaces/cli/runner.go`) under a unique name.  
3. Users can then pass `--exchange your_adapter_name` to target that exchange.

---

## License

Copyright 2025 Isak Ruas

Licensed under the Apache License, Version 2.0 (the 'License');
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an 'AS IS' BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
