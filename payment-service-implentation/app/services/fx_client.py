import requests
from requests.exceptions import RequestException, Timeout
import logging
import time

logger = logging.getLogger(__name__)

class FXServiceError(Exception):
    """Custom exception to standardize errors coming from the FX service."""
    pass

class FXClient:
    def __init__(self, base_url: str = "http://localhost:4000", max_retries: int = 3):
        self.quote_url = f"{base_url}/twirp/payments.v1.FXService/GetQuote"
        self.timeout = 5.0 
        self.max_retries = max_retries

    def get_exchange_rate(self, source: str, target: str) -> float:
        """
        Fetches the exchange rate for a currency pair with exponential backoff retries.
        Raises FXServiceError if the rate cannot be obtained after all attempts.
        """
        payload = {
            "source_currency": source,
            "target_currency": target
        }
        
        for attempt in range(1, self.max_retries + 1):
            try:
                response = requests.post(self.quote_url, json=payload, timeout=self.timeout)
                response.raise_for_status() 
                
                data = response.json()
                rate = data.get("exchange_rate")
                
                if rate is None or not isinstance(rate, (int, float)) or rate <= 0:
                    logger.error(f"Invalid rate received: {rate}")
                    raise FXServiceError("FX Service returned an invalid or missing exchange rate.")
                    
                return float(rate)
                
            except Timeout:
                logger.warning(f"Attempt {attempt}/{self.max_retries} - Timeout connecting to FX Service.")
                if attempt == self.max_retries:
                    logger.error("Max retries reached. FX Service is unavailable.")
                    raise FXServiceError(f"FX Service timed out after {self.max_retries} attempts.")
                    
                time.sleep(2 ** attempt)
                
            except RequestException as e:
                logger.warning(f"Attempt {attempt}/{self.max_retries} - Network error: {e}")
                if attempt == self.max_retries:
                    raise FXServiceError(f"Network error after {self.max_retries} attempts: {str(e)}")
                    
                time.sleep(2 ** attempt)
                
            except ValueError:
                logger.error("Failed to parse JSON from FX Service.")
                raise FXServiceError("FX Service returned malformed JSON.")