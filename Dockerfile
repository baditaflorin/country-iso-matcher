# Dockerfile
# ---- Builder Stage ----
FROM golang:1.23-alpine AS builder

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
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /server /app/server

# Copy data files (CSV/TSV)
COPY --from=builder /app/data /app/data

# Copy web GUI files
COPY --from=builder /app/web /app/web

# Copy example configuration
COPY --from=builder /app/configs /app/configs

# Expose service port
EXPOSE 3030

# Set environment variables for default CSV data source
ENV DATA_SOURCE=csv
ENV DATA_COUNTRIES_FILE=data/countries.csv
ENV DATA_ALIASES_FILE=data/aliases.csv
ENV GUI_ENABLED=true
ENV GUI_PATH=/admin

# Run app
ENTRYPOINT ["/app/server"]
