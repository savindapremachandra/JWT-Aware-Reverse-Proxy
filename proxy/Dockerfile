FROM golang:1.24.3 AS builder
WORKDIR /app

# Copy go.mod and go.sum first (for layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the binary
RUN go build -o proxy .

# Final image
FROM debian:bookworm-slim
WORKDIR /app

# Copy binary and public.pem
COPY --from=builder /app/proxy .
COPY public.pem .
COPY routes.yaml .

EXPOSE 9000
CMD ["./proxy"]
