from pydantic import BaseModel, Field
from typing import Optional
from enum import Enum
from datetime import datetime, timezone
import uuid

class PaymentStatus(str, Enum):
    PENDING = "pending"
    SUCCEEDED = "succeeded"  
    FAILED = "failed"        

class PaymentRequest(BaseModel):
    """Data model representing the incoming payload from the user."""
    sender: str              
    receiver: str            
    # Rationale: Field(gt=0) enforces the rule that inputs must be valid.
    amount: float = Field(..., gt=0, description="Payment amount must be strictly positive")
    source_currency: str = Field(..., min_length=3, max_length=3, description="e.g., USD")
    destination_currency: str = Field(..., min_length=3, max_length=3, description="e.g., EUR")

class PaymentRecord(BaseModel):
    """Data model representing the stored state of the payment."""
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    request: PaymentRequest
    status: PaymentStatus = PaymentStatus.PENDING

    # These fields are Optional because they cannot be populated 
    # until the FX Service is successfully called.
    payout_amount: Optional[float] = None  
    exchange_rate: Optional[float] = None
    
    # Captures errors for later retrieval
    diagnostic_info: Optional[str] = None
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))