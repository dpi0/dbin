# Build stage
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOARCH=amd64 GOOS=linux go build -o ./bin/dbin ./cmd/dbin

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/dbin .

COPY --from=builder /app/web ./web/

EXPOSE 1323

CMD ["./dbin"]