# Dockerfile
# ---- Builder Stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy go.mod and go.sum first (to leverage Docker layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application (note: using src/cmd/server path)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server ./src/cmd/server

# ---- Final Stage ----
FROM scratch

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy static binary from builder
COPY --from=builder /server /server

# Expose service port
EXPOSE 3030

# Run app
ENTRYPOINT ["/server"]