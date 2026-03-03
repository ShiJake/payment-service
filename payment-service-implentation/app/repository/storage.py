from typing import Dict, Optional
from app.models.payment import PaymentRecord

class PaymentRepository:
    """
    An in-memory storage repository for PaymentRecord objects.
    """
    def __init__(self):
        self._storage: Dict[str, PaymentRecord] = {}

    def save(self, record: PaymentRecord) -> PaymentRecord:
        """
        Stores or updates a payment record in the repository.
        """
        self._storage[record.id] = record
        return record

    def get_by_id(self, payment_id: str) -> Optional[PaymentRecord]:
        """
        Retrieves a payment record by its unique identifier.
        Returns None if the record does not exist.
        """
        return self._storage.get(payment_id)

# Global repository instance to be used across the application
payment_db = PaymentRepository()