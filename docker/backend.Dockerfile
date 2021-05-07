FROM golang:1.16-alpine AS builder

WORKDIR /usr/src/app

RUN apk --no-cache add libgit2-dev gcc musl-dev

COPY . .

RUN go build -o haas .

FROM alpine:latest

WORKDIR /usr/src/app

RUN apk --no-cache add libgit2

COPY --from=builder /usr/src/app/haas ./haas

CMD ["./haas"]
