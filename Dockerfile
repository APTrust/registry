# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

RUN apk add git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /main

EXPOSE 80

CMD [ "/main" ]
