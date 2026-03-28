FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1001

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate templ files
RUN templ generate

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/web ./cmd/web

# Runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/web .

EXPOSE 8084

CMD ["./web"]
