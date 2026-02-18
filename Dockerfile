FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY docs/ ./docs

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/file-parser ./cmd/api/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/bin/file-parser ./

CMD ["./file-parser"]
