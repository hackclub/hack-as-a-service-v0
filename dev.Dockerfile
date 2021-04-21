FROM golang:1.16-alpine

COPY . /code

WORKDIR /code

RUN go get github.com/cespare/reflex

CMD go run -v .
