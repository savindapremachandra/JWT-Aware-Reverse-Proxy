# decode_jwt.py
import jwt
from jwt import InvalidTokenError
from datetime import timezone

with open("pki/public.pem", "r") as f:
    public_key = f.read()

# Replace with your token string
token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiYWRhbSIsInRpZXIiOiJzdWJzY3JpYmVkIiwiZXhwIjoxNzQ4NTM1ODcxfQ.vxJNgkQolnlYQi3lFyC-EQrJ1HXJ7CHGJXqaxFsAPU8cP4Cfb7YHqsdCT7iltvf89QRy-B-SvQh4dV7-V__ReZaHubo-NVa2ugeDK1_d63xApeJ2dJrwaCdqO0gm-Jbv3uvzAM39vCzQtSsvtZO_cmFPxZr2Uws29h2g_P1fUxlbxG2laXsCGhARDXuNBSU0HvYVx8HatVhD7seuP6TgML4Z-Z8HIraOnt5cyjLUNxLhxHGPCi6aW84u6aR1wGFkYHcqWPxNeqwBAS_-NxyuTUo6O0T1wXvitHFzW36rtguhmLa8ZRS1GaoKJvZhv3VufjQQ4NOesnAsk-yaYXPLDw"

try:
    decoded = jwt.decode(token, public_key, algorithms=["RS256"])
    print("✅ Token is valid!")
    print("Payload:", decoded)
except InvalidTokenError as e:
    print("❌ Invalid token:", e)
