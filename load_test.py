import urllib.request
import urllib.error
import json
import concurrent.futures
import threading
import random
import time

# Configuration
BASE_URL = "http://localhost:8080"
STOCKS = ["AAPL", "MSFT", "GOOGL", "TSLA"]
WALLETS = ["wallet1", "wallet2", "wallet3", "wallet4", "wallet5"]

THREADS = 100
REQUESTS_PER_THREAD = 100
CHAOS_DELAY = 1.5


def do_request(method, path, data=None):
    url = BASE_URL + path
    headers = {'Content-Type': 'application/json'}
    body = json.dumps(data).encode('utf-8') if data else None

    req = urllib.request.Request(url, data=body, headers=headers, method=method)

    try:
        with urllib.request.urlopen(req) as response:
            return response.status, response.read()
    except urllib.error.HTTPError as e:
        return e.code, e.read()
    except urllib.error.URLError as e:
        return 0, str(e.reason)


def setup_bank():
    print("Step 1: Initializing bank state (loading stocks)")

    payload = {
        "stocks": [{"name": s, "quantity": 10000} for s in STOCKS]
    }

    status, _ = do_request("POST", "/stocks", payload)

    if status == 200:
        print("Bank initialized successfully (10,000 units per stock).")
    else:
        print(f"Failed to initialize bank. Status: {status}")


def trigger_chaos(delay):
    time.sleep(delay)

    print(f"\n[CHAOS TEST] Triggering failure via {BASE_URL}/chaos")

    status, _ = do_request("POST", "/chaos")

    print(f"[CHAOS TEST] Chaos triggered. Response status: {status}\n")


def worker(worker_id):
    success = 0
    client_errors = 0
    server_errors = 0

    for _ in range(REQUESTS_PER_THREAD):
        wallet = random.choice(WALLETS)
        stock = random.choice(STOCKS)
        action = random.choice(["buy", "sell"])

        status, _ = do_request(
            "POST",
            f"/wallets/{wallet}/stocks/{stock}",
            {"type": action}
        )

        if status == 200:
            success += 1
        elif status in (400, 404):
            client_errors += 1
        else:
            server_errors += 1

    return success, client_errors, server_errors


def run_load_test():
    print("\nStep 2: Starting load test")
    print(f"Threads: {THREADS}, Requests per thread: {REQUESTS_PER_THREAD}")
    print(f"Total requests: {THREADS * REQUESTS_PER_THREAD}")

    start_time = time.time()

    total_success = 0
    total_client_errors = 0
    total_server_errors = 0

    chaos_thread = threading.Thread(target=trigger_chaos, args=(CHAOS_DELAY,))
    chaos_thread.start()

    with concurrent.futures.ThreadPoolExecutor(max_workers=THREADS) as executor:
        futures = [executor.submit(worker, i) for i in range(THREADS)]

        for future in concurrent.futures.as_completed(futures):
            s, c_err, s_err = future.result()
            total_success += s
            total_client_errors += c_err
            total_server_errors += s_err

    duration = time.time() - start_time

    print("\nTest completed")
    print(f"Duration: {duration:.2f} seconds")
    print(f"Successful requests (200): {total_success}")
    print(f"Client errors (400/404): {total_client_errors}")
    print(f"Server errors (500/502/503): {total_server_errors}")
    print(f"Throughput: {(THREADS * REQUESTS_PER_THREAD) / duration:.2f} req/s")


def fetch_logs():
    print("\nStep 3: Fetching audit logs")

    time.sleep(1)

    status, body = do_request("GET", "/log")

    if status == 200:
        data = json.loads(body)
        logs = data.get("log", [])

        print(f"Retrieved {len(logs)} audit log entries")

        with open("audit_logs.json", "w", encoding="utf-8") as f:
            json.dump(logs, f, indent=4)

        print("Audit logs saved to audit_logs.json")
    else:
        print(f"Failed to fetch logs. Status: {status}")


if __name__ == "__main__":
    setup_bank()
    run_load_test()
    fetch_logs()
