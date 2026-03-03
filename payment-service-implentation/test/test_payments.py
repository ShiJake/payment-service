import pytest
from fastapi.testclient import TestClient
from unittest.mock import patch, MagicMock
from requests.exceptions import Timeout

from app.main import app
from app.models.payment import PaymentStatus
from app.services.fx_client import FXClient, FXServiceError

# Send HTTP requests directly to our FastAPI app.
client = TestClient(app)

def test_create_payment_success():
    """Tests the ideal path: Valid inputs and a successful FX rate retrieval."""
    payload = {
        "sender": "Alice",
        "receiver": "Bob",
        "amount": 100.0,
        "source_currency": "USD",
        "destination_currency": "EUR"
    }

    with patch("app.services.payment_manager.fx_client.get_exchange_rate", return_value=0.92):
        response = client.post("/payments", json=payload)
        assert response.status_code == 201
        data = response.json()
        assert data["status"] == PaymentStatus.SUCCEEDED
        assert data["payout_amount"] == 92.0

def test_create_payment_invalid_input():
    """Tests the fail-fast constraint: Invalid inputs should be rejected immediately."""
    payload = {
        "sender": "Alice",
        "receiver": "Bob",
        "amount": -50.0,  
        "source_currency": "US",  
        "destination_currency": "EUR"
    }
    response = client.post("/payments", json=payload)
    assert response.status_code == 422

def test_create_payment_fx_failure():
    """Tests resilience: If the FX service fails, the payment must fail gracefully."""
    payload = {
        "sender": "Alice",
        "receiver": "Bob",
        "amount": 100.0,
        "source_currency": "USD",
        "destination_currency": "EUR"
    }

    mock_error_message = "FX Service timed out after 3 attempts."
    with patch("app.services.payment_manager.fx_client.get_exchange_rate", side_effect=FXServiceError(mock_error_message)):
        response = client.post("/payments", json=payload)
        assert response.status_code == 201  
        data = response.json()
        assert data["status"] == PaymentStatus.FAILED
        assert data["diagnostic_info"] == mock_error_message

def test_get_payment_success():
    """
    Tests the retrieval endpoint to ensure stored payment details can be fetched.
    """
    payload = {
        "sender": "Charlie",
        "receiver": "Diana",
        "amount": 50.0,
        "source_currency": "GBP",
        "destination_currency": "JPY"
    }

    # 1. Create a payment first
    with patch("app.services.payment_manager.fx_client.get_exchange_rate", return_value=150.0):
        post_response = client.post("/payments", json=payload)
        payment_id = post_response.json()["id"]

    # 2. Retrieve the payment using the generated ID
    get_response = client.get(f"/payments/{payment_id}")
    
    assert get_response.status_code == 200
    data = get_response.json()
    assert data["id"] == payment_id
    # Proves the stored state is accurate (50.0 * 150.0 = 7500.0)
    assert data["payout_amount"] == 7500.0
    assert data["status"] == PaymentStatus.SUCCEEDED

def test_get_payment_not_found():
    """
    Tests that requesting a non-existent payment returns a standard 404 response.
    """
    # Query an ID that we know is not in the repository
    response = client.get("/payments/invalid-uuid-1234")
    assert response.status_code == 404
    assert response.json()["detail"] == "Payment with ID invalid-uuid-1234 not found."

@patch("app.services.fx_client.time.sleep")
@patch("app.services.fx_client.requests.post")
def test_fx_client_retries_and_succeeds(mock_post, mock_sleep):
    """
    Tests the Exponential Backoff logic inside the FXClient directly.
    """
    fx = FXClient(max_retries=3)
    
    # Create a mock successful response for the final attempt
    mock_success_response = MagicMock()
    mock_success_response.json.return_value = {"exchange_rate": 0.85}
    
    # Configure requests.post to fail with a Timeout TWICE, then succeed ONCE.
    mock_post.side_effect = [
        Timeout("Simulated Timeout 1"), 
        Timeout("Simulated Timeout 2"), 
        mock_success_response
    ]
    
    rate = fx.get_exchange_rate("USD", "GBP")
    
    assert rate == 0.85
    # Client attempted the network call exactly 3 times
    assert mock_post.call_count == 3
    # Client triggered the time.sleep() backoff exactly twice before succeeding
    assert mock_sleep.call_count == 2