FROM golang:stretch

COPY . /code

WORKDIR /code

RUN go get github.com/cespare/reflex

CMD go run -v .
