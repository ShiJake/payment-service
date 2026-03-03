# Cross-Currency Payment Service

This repository contains the implementation of a cross-currency payment service. It is designed to accept payment requests, fetch live exchange rates from an external (and potentially volatile) FX Service, compute the final payout, and reliably store the transaction state.

## Architectural Design & Tradeoffs
To ensure thoughtful design, this application utilizes a **Layered Architecture**:
* **Domain Models (Pydantic):** Strictly enforces data validation at the application boundary. It mathematically guarantees that currencies are 3-character ISO codes and amounts are strictly positive before any processing begins.
* **Service Layer:** Orchestrates business logic independently of the web framework. The `FXClient` includes explicit resilience patterns (5-second timeouts and HTTP error trapping) to gracefully handle the FX service's simulated network latency and intermittent failures.
* **Repository Pattern:** Abstracts data persistence. An in-memory dictionary is used for the scope of this assessment, allowing for $O(1)$ read/write speeds, but it can be easily swapped for a persistent SQL database in a production environment.

---

## Prerequisites
* **Node.js & npm** (To run the provided FX Service)
* **Python 3.13+** (To run the Payment API)

---

## Getting Started

You must run both the external FX Service and the Payment Service concurrently. Open two separate terminal windows.

### 1. Start the FX Service
Navigate to the FX service directory and start the Node application:
```bash
cd fx-service
npm install
npm start
```
Note: The FX service simulates real-world conditions and runs on `http://localhost:4000`.

### 2. Start the Payment Service
Navigate to the new payment implementation directory, install the dependencies, and start the FastAPI server:
```bash
cd payment-service-implentation
pip install -r requirements.txt
python -m uvicorn app.main:app --reload --port 8000
```
Note: The Payment API will run on `http://localhost:8000`.

### 3. Running the Tests
To verify the system's logic and error handling without relying on the external FX service, 
execute the unit test suite from the `payment-service-implentation` directory:
```bash
python -m pytest test/
```

## API Usage Examples
Below are examples of how to interact with the service using PowerShell (Windows native) and standard cURL (Unix/Linux).

### Process a Payment 

Submits a payment request for processing.

PowerShell (Windows):
```powershell
$body = @{
    sender = "Alice"
    receiver = "Bob"
    amount = 100.0
    source_currency = "USD"
    destination_currency = "EUR"
} | ConvertTo-Json

Invoke-RestMethod -Method Post -Uri "http://localhost:8000/payments" -ContentType "application/json" -Body $body
```

cURL (Unix/Linux/Git Bash):
```bash
curl -X POST http://localhost:8000/payments \
  -H "Content-Type: application/json" \
  -d '{
    "sender": "Alice",
    "receiver": "Bob",
    "amount": 100.0,
    "source_currency": "USD",
    "destination_currency": "EUR"
  }'
```

### Retrieve a Payment
Retrieves the payment status, calculated payout amount, and any diagnostic information. Replace `{id}` with the UUID returned from the POST request.

PowerShell (Windows):
```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8000/payments/{id}"
```

cURL (Unix/Linux/Git Bash):
```bash
curl http://localhost:8000/payments/{id}
```