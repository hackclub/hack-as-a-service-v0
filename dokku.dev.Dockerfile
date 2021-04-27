FROM --platform=amd64 golang:stretch AS builder

COPY . /code

WORKDIR /code/dokkud

RUN go build .

FROM --platform=amd64 golang:stretch AS plugin_builder

WORKDIR /code

COPY . .

RUN ./install

FROM dokku/dokku:latest

COPY --from=plugin_builder /code /var/lib/dokku/plugins/available/haas
RUN ln -s /var/lib/dokku/plugins/available/haas /var/lib/dokku/plugins/enabled/haas

COPY --from=builder /code/dokkud/dokkud /opt/dokkud
COPY ./dokkud/20_dokkud_init /etc/my_init.d/
RUN chmod +x /etc/my_init.d/20_dokkud_init
