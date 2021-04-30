## BUILD FRONTEND ##
FROM node:14-alpine AS frontend

WORKDIR /usr/src/app

COPY frontend/ ./frontend
COPY ./swagger.yaml ./swagger.yaml

WORKDIR /usr/src/app/frontend

RUN yarn install

RUN yarn build
## BUILD FRONTEND ##

## BUILD BACKEND ##
FROM golang:1.16-alpine AS backend

WORKDIR /usr/src/app

COPY . .
COPY --from=frontend /usr/src/app/frontend/out/ ./frontend/out

RUN go build -o haas .
## BUILD BACKEND ##

## RUN APP ##
FROM alpine:latest

WORKDIR /usr/src/app

COPY --from=backend /usr/src/app/haas/ ./haas
COPY --from=backend /usr/src/app/swagger.yaml ./swagger.yaml
COPY --from=frontend /usr/src/app/frontend/out/ ./frontend/out

CMD ["/usr/src/app/haas"]
## RUN APP ##