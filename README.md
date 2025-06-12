# 🛡️ JWT-Aware Reverse Proxy in Go

This is a proof-of-concept reverse proxy written in Go 🐹 that routes traffic based on the `tier` claim in a JWT token. It's built to help understand how JWT authentication and reverse proxying can work together — especially in tiered-access API systems.

---

## 🧠 What It Does

- Verifies JWT tokens using RSA public/private key cryptography
- Extracts the `tier` claim (e.g., `free`, `subscribed`)
- Routes requests to different backend services based on the tier
- Everything runs inside containers on Minikube for easy testing

---

## 🏗️ Architecture

```bash
[ Client ] --> [ Reverse Proxy (Go) ] --> [ Free Tier API ]
                                   |--> [ Subscribed Tier API ]
````

* JWT is passed via `Authorization: Bearer <token>`
* Proxy reads routing logic from `config.yaml`
* Backends are basic Go servers that respond based on the tier

---

## 📦 Components

### ✅ Reverse Proxy

* Uses `httputil.NewSingleHostReverseProxy`
* JWT verification with RSA public key (`public.pem`)
* Routing logic based on `tier` claim in the JWT
* Example config:

  ```yaml
  routes:
    free: http://free-tier-service
    subscribed: http://subscribed-tier-service
  ```

### 🔐 JWT Token

* Signed with a private key (`private.pem`)
* Example claim:

  ```json
  {
    "sub": "user@example.com",
    "tier": "subscribed",
    "exp": 9999999999
  }
  ```

---

## 🚀 Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/savindapremachandra/jwt-reverse-proxy-go.git
cd jwt-reverse-proxy-go
```

### 2. Generate RSA Keys

```bash
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem
```

### 3. Build Docker Images

```bash
docker build -t jwt-proxy ./proxy
docker build -t free-api ./free-api
docker build -t sub-api ./subscribed-api
```

---

## 🧪 Testing

Send requests like this (with valid JWT):

```bash
curl -H "Authorization: Bearer <token>" http://<proxy-service>:<port>/
```

Change the `tier` claim in your JWT to test different routing behavior.

---

## 🌍 Real-World Use Cases

* Tier-based API access (e.g., free vs premium)
* Feature gating in SaaS platforms
* Content delivery services with subscription levels
* Internal APIs with partner-level segmentation
* Educational content portals with enrollment checks

---

## ⚠️ Not for Production (Yet)

This is a **test project**, and while it's fun and useful to learn from, it comes with a few caveats:

* No HTTPS
* No token revocation or refresh flow
* No backend validation of JWT
* No observability, metrics, or rate limiting
* No public key hot reload or rotation

For real-world usage, consider adding those features or using a dedicated API gateway.

---

## 📚 Blog Post

Check out the full explanation and code breakdown in this [blog post](#) 📝

---

## 🤝 Contributions

Got ideas, questions, or feedback? Feel free to open issues or PRs — or just drop a star ⭐ if you liked it!

---

## 📄 License

MIT. Use it, learn from it, build something better with it.

```

---

Let me know if you want a badge-style header (build, version, etc.), a Makefile, or instructions for generating sample JWTs using `jwt.io` or Go scripts 👨‍💻
```
