## BUILD FRONTEND ##
FROM node:14-alpine

WORKDIR /usr/src/app

COPY frontend/ .

RUN yarn install

RUN yarn build
## BUILD FRONTEND ##

## BUILD BACKEND ##
FROM golang:1.16-alpine

WORKDIR /usr/src/app

COPY . .
COPY --from=0 /usr/src/app/out/ ./frontend/out

RUN go build -o haas .
## BUILD BACKEND ##

## RUN APP ##
FROM alpine:latest

WORKDIR /usr/src/app

COPY --from=1 /usr/src/app/haas/ ./haas
COPY --from=0 /usr/src/app/out/ ./frontend/out

CMD ["/usr/src/app/haas"]
## RUN APP ##