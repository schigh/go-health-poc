FROM golang:1.9-alpine
MAINTAINER Steve High <steve.high@gmail.com>
COPY ./ /go/src/github.com/schigh/go-health-poc
WORKDIR /go/src/github.com/schigh/go-health-poc

RUN apk add --no-cache git bash
RUN go get -u goji.io
RUN go get github.com/go-redis/redis
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/sirupsen/logrus
RUN mv go-health ../
RUN go build -o app main.go

EXPOSE 80
