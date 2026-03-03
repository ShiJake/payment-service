from app.models.payment import PaymentRequest, PaymentRecord, PaymentStatus
from app.repository.storage import payment_db
from app.services.fx_client import FXClient, FXServiceError

# Instantiate the client once to be reused across all payment requests.
fx_client = FXClient()

def process_payment(request: PaymentRequest) -> PaymentRecord:
    """
    Orchestrates the payment lifecycle: creation, FX rate retrieval, 
    payout calculation, and status updates.
    """
    # Initialize and store the pending payment 
    # Save immediately so we have a record even if the application crashes during the FX call.
    record = PaymentRecord(request=request)
    payment_db.save(record)
    
    try:
        # Call the external FX Service 
        rate = fx_client.get_exchange_rate(
            source=request.source_currency, 
            target=request.destination_currency
        )
        
        # Compute the payout and update status on success
        record.exchange_rate = rate
        # The payment is computed by multiplying the amount by the FX rate.
        record.payout_amount = request.amount * rate
        record.status = PaymentStatus.SUCCEEDED
        
    except FXServiceError as e:
        # Handle failures gracefully and store diagnostic info
        # If the FX service is slow or unavailable, the payment gracefully fails 
        # and we capture the error for later retrieval.
        record.status = PaymentStatus.FAILED
        record.diagnostic_info = str(e)
        
    # Save the final state back to the repository
    return payment_db.save(record)