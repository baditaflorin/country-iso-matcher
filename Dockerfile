# ---- Builder stage ----
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/country-iso-matcher ./src/cmd/server

# ---- Final stage ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget tini \
 && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /out/country-iso-matcher /app/country-iso-matcher
COPY --from=builder /app/data     /app/data
COPY --from=builder /app/configs  /app/configs
COPY --from=builder /app/web      /app/web
COPY --from=builder /app/service.yaml /app/service.yaml
USER app

ENV PORT=18315 \
    SERVER_HOST=0.0.0.0 \
    SERVER_ENVIRONMENT=production \
    DATA_SOURCE=csv \
    DATA_COUNTRIES_FILE=data/countries.csv \
    DATA_ALIASES_FILE=data/aliases.csv \
    LOG_LEVEL=info \
    LOG_FORMAT=json \
    GUI_ENABLED=false

EXPOSE 18315
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1:${PORT:-18315}/health >/dev/null || exit 1
ENTRYPOINT ["/sbin/tini","--"]
CMD ["/app/country-iso-matcher"]
