FROM --platform=amd64 golang:stretch AS plugin_builder

WORKDIR /code

COPY . .

RUN ./install

FROM dokku/dokku:latest

WORKDIR /tmp

RUN curl -L -o golang.tar.gz https://dl.google.com/go/go1.16.3.linux-amd64.tar.gz && \
    tar -xvf golang.tar.gz && \
    mv go /opt && \
    mkdir -p /opt/gopath && mkdir -p /opt/gocache
ENV PATH="/opt/go/bin:/opt/gopath/bin:${PATH}" GOROOT="/opt/go" GOPATH="/opt/gopath" GOCACHE="/opt/gocache"

RUN go get github.com/cespare/reflex

COPY --from=plugin_builder /code /var/lib/dokku/plugins/available/haas
RUN ln -s /var/lib/dokku/plugins/available/haas /var/lib/dokku/plugins/enabled/haas

COPY ./dokkud/20_dokkud_init /etc/my_init.d/
RUN chmod +x /etc/my_init.d/20_dokkud_init
