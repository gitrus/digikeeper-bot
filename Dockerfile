FROM golang:1.24.0-bookworm AS builder

WORKDIR /usr/app

COPY go.mod go.sum ./

RUN go mod download && mkdir -p /usr/app/bin

COPY . .

RUN go build -o /usr/app/bin/digikeeper-bot ./cmd/bot

# main
FROM debian:bookworm-slim AS main

COPY --from=builder /usr/app/bin/digikeeper-bot /usr/app/bin/digikeeper-bot

CMD ["/usr/app/bin/digikeeper-bot"]

EXPOSE 8081
# metrics
EXPOSE 8091
