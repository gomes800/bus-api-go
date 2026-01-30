# ---- build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates tzdata

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/app ./cmd/api

# ---- run stage ----
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H appuser
USER appuser
COPY --from=builder /bin/app /app/app
EXPOSE 8080
ENTRYPOINT ["/app/app"]
