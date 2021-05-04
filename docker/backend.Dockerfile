FROM golang:1.16-alpine AS builder

WORKDIR /usr/src/app

COPY . .

RUN go build -o haas .

FROM alpine:latest

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/haas ./haas

CMD ["./haas"]
