# Stock Market Simulator

Simple REST service simulating a basic stock market with wallets and a central bank. Designed with High Availability (HA) and concurrent transaction safety in mind.

## Tech stack
* **Go** (Gin framework) - Lightweight and fast.
* **Redis** (Database & Lua scripts) - Ensures atomicity and prevents race conditions.
* **Docker & Docker Compose** - Containerization for multi-arch support.
* **Nginx** - Acts as a Load Balancer ensuring High Availability.

## Features
* Buy / Sell stocks (Atomic lock-free operations via Redis Lua scripts).
* Wallet state management & Bank as a single liquidity provider.
* High Availability (Multiple Go backend instances behind Nginx).
* Graceful shutdown (Handles SIGTERM/SIGINT properly).
* Multi-architecture support out of the box (x64 / arm64).

## Requirements
* Docker & Docker Compose
* Python 3 (optional, for load/chaos testing)
* Go (optional, for running local E2E tests, everything else is done in containers)

## Run

The solution can be started using a single command where the parameter is the desired port.

**On Linux / macOS:**
```bash
chmod +x run.sh
./run.sh 8080
```
**On Windows**
```cmd
run.bat 8080
```
The application will be available at: http://localhost:8080 (or whichever port you specified).
## API Endpoints

### Wallet Operations

Buy / Sell stock:

```
POST /wallets/{wallet_id}/stocks/{stock_name}
Body: {"type": "buy"} or {"type": "sell"}
```

Get wallet state:

```
GET /wallets/{wallet_id}
```

Get single stock qty:
```
GET /wallets/{wallet_id}/stocks/{stock_name}
```
### Bank Operations

Set bank state:
```
POST /stocks

Body: {"stocks":[{"name":"AAPL", "quantity":99}]}
```
Get bank state:
```
GET /stocks
```
### System & Audit

Get audit log:
```
GET /log (Returns up to 10,000 successful operations in order of occurrence)
```
Health check:
```
GET /health
```
Chaos (kill instance):
```
POST /chaos
```
## Testing

* End-to-End Tests (Go)
Simulates the entire market flow using an in-memory Redis instance (miniredis):
```go
go test -v ./app/tests/
```
* Load & Chaos Test (Python)
Spawns multiple threads to simulate heavy traffic and injects a node failure (/chaos) to test HA resilience:
```python
python3 load_test.py
```

## Notes

### Thread Safety
Concurrency and race conditions are fully mitigated using Redis Lua scripts. Even with hundreds of concurrent requests, it is impossible to buy more stocks than the bank possesses.

### High Availability
Killing one backend instance via /chaos does not drop the entire system. Nginx transparently routes incoming traffic to the surviving instance and tries to restart the closed one.

### Chaos Endpoint Behavior
Invoking /chaos deliberately forces os.Exit(1) without returning an HTTP response. By design, Nginx is configured to return 502 Bad Gateway to the caller of this specific endpoint and is restricted from retrying it on the healthy node (proxy_next_upstream off;).

### Validation & Errors
Input data is validated using Gin binding tags. The application leverages Go's built-in errors.Is for clean error handling.
