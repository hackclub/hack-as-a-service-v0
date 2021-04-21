FROM golang:stretch AS builder

COPY . /code

WORKDIR /code/dokkud

RUN go build .

FROM dokku/dokku:latest

COPY --from=builder /code/dokkud/dokkud /opt/dokkud
COPY ./dokkud/20_dokkud_init /etc/my_init.d/
RUN chmod +x /etc/my_init.d/20_dokkud_init
