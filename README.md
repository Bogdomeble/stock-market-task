# Stock Market Simulator

Simple REST service simulating a basic stock market with wallets and a central bank. Designed with High Availability (HA) and concurrent transaction safety in mind.

## Tech stack
* Go (Gin framework)
* Redis (Database & Lua scripts for atomicity)
* Docker & Docker Compose
* Nginx (Load balancing / HA)

## Features
* Buy / Sell stocks (Atomic operations via Redis Lua scripts)
* Wallet state management & Bank as single liquidity provider
* High availability (Multiple Go instances behind Nginx)
* Graceful shutdown (Handles SIGTERM/SIGINT)
* Chaos testing ready (Load balancer handles failovers transparently)

## Requirements
* Docker & Docker Compose
* Python 3 (optional, for load/chaos testing)
* Go (optional, for running local E2E tests)

## Run

Solution can be started using a single command where the parameter is the desired port:

```bash

PORT=8080 docker compose up --build
```

```powershell

(Windows PowerShell: $env:PORT="8080"; docker compose up --build)

```
Application will be available at: http://localhost:8080 
## Testing

1. End-to-End Tests (Go)
Simulates the entire market flow using an in-memory Redis instance (miniredis):

```go
go test -v ./app/tests/
```

2. Load & Chaos Test (Python)
Spawns multiple threads to simulate heavy traffic and injects a node failure (/chaos) to test HA resilience:

```python
python3 load_test.py
```

## API

### Wallet Operations

```
Buy / Sell stock: POST /wallets/{wallet_id}/stocks/{stock_name} | Body: {"type": "buy" | "sell"}
```

```
Get wallet state: GET /wallets/{wallet_id}
```

```
Get single stock qty: GET /wallets/{wallet_id}/stocks/{stock_name}
```

### Bank Operations 

```
Set bank state: POST /stocks | Body: {"stocks":[{"name":"AAPL", "quantity":99}]}
``` 

```
Get bank state: GET /stocks
```

### System & Audit

```
Get audit log: GET /log (Returns up to 10,000 successful operations)
```

```
Health check: GET /health
```

```
Chaos (kill instance): POST /chaos
```
### Notes

Stock price is fixed to 1. No wallet balance or funds tracking.

Performance: Each Go backend consumes only ~10-15MB RAM and handles >1500 RPS locally.

Thread Safety: Concurrency and race conditions are fully mitigated using Redis Lua scripts.

High Availability: Killing one backend instance via /chaos does not drop active connections; Nginx transparently routes traffic to the surviving instance.