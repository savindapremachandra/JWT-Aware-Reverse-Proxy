Building a JWT-Aware Reverse Proxy in Go for Tiered API Access

Introduction

In this blog post, I’ll walk you through a small proof-of-concept project I built: a reverse proxy in Go that routes requests to different backend services based on the tier claim in a JWT token.

The proxy uses RSA public/private key cryptography to validate the JWT, then forwards the request to either a free-tier or subscribed-tier backend. The whole setup runs in containers on Minikube. It’s a great way to understand how reverse proxies and JWT authentication can work together — even if it’s not production-grade just yet.

Motivation

I wanted to understand how reverse proxies work at a lower level and how JWTs could be used to manage access control between service tiers. I also wanted something I could deploy easily with Docker and Minikube, so it’s very infrastructure- and developer-friendly.

Architecture Overview

The client sends a request with a JWT in the Authorization header
The reverse proxy verifies the JWT using a public RSA key
It extracts the tier claim (e.g., "free" or "subscribed")
Based on that, it routes the request to the appropriate backend service

All services are containerized and deployed in Minikube.

