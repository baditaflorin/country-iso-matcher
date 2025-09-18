# ---- Builder Stage ----
# Use the official Go image to build the application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.* ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application, creating a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server .

# ---- Final Stage ----
# Use a minimal 'scratch' image which contains nothing but our binary
FROM scratch

WORKDIR /

# Copy the static binary from the builder stage
COPY --from=builder /server .

# Expose the port the app runs on
EXPOSE 3030

# Command to run the application
ENTRYPOINT ["/server"]