# ---- Builder Stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first (to leverage Docker layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY ./src ./src

# Build the Go application (point at cmd package)
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./src/cmd

# ---- Final Stage ----
FROM scratch

WORKDIR /

# Copy static binary from builder
COPY --from=builder /server .

# Expose service port
EXPOSE 3030

# Run app
ENTRYPOINT ["/server"]
