from fastapi import FastAPI, HTTPException, status
from app.models.payment import PaymentRequest, PaymentRecord
from app.services.payment_manager import process_payment
from app.repository.storage import payment_db

# Initialize the FastAPI application
app = FastAPI(title="Cross-Currency Payment Service")

@app.post("/payments", response_model=PaymentRecord, status_code=status.HTTP_201_CREATED)
def create_payment(request: PaymentRequest):
    """
    Receives a payment request, processes it via the FX service, 
    and returns the resulting payment record.
    """
    # FastAPI automatically validates the incoming JSON against 
    # the PaymentRequest Pydantic model before this line is even reached.
    # If the inputs are invalid, it automatically returns a 422 Unprocessable Entity error.
    record = process_payment(request)
    return record

@app.get("/payments/{payment_id}", response_model=PaymentRecord)
def get_payment(payment_id: str):
    """
    Retrieves the details of a specific payment by its ID.
    """
    # Query the in-memory storage repository abstractly.
    record = payment_db.get_by_id(payment_id)
    
    if not record:
        # Returning a strict 404 is standard RESTful practice 
        # when a requested resource does not exist.
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, 
            detail=f"Payment with ID {payment_id} not found."
        )
        
    # By returning the PaymentRecord, we automatically expose the payment status,
    # payout amount, and diagnostic information.
    return record