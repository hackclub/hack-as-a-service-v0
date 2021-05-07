FROM golang:stretch

RUN echo "deb http://deb.debian.org/debian stretch-backports main" >> /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y libgit2-dev/stretch-backports && \
    rm -rf /var/lib/apt/lists/*

COPY . /code

WORKDIR /code

RUN go get .

RUN go get github.com/cespare/reflex

CMD go run -v .
