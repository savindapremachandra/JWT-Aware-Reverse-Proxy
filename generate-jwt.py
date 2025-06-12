# token_gen.py
import jwt, time
from datetime import datetime, timedelta, timezone

with open("pki/private-key.pem", "r") as f:
    private_key = f.read()

payload = {
    "user": "adam",
    "tier": "subscribed",  # or "free"
    "exp": datetime.now(timezone.utc) + timedelta(hours=1)
}

token = jwt.encode(payload, private_key, algorithm="RS256")
print(token)