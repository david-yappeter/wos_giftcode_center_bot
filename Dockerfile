FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN apk add --no-cache build-base \
    && go mod download

COPY . .

RUN CGO_ENABLED=1 go build -tags http -o /app/build/wos .

FROM alpine:latest

COPY --from=builder /app/build/wos /app/build/wos

WORKDIR /app/build/

CMD ["./wos"]
