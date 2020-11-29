FROM golang:1.13-alpine

ARG MYSQL_HOST
ARG MYSQL_PORT
ARG MYSQL_DATABASE
ARG MYSQL_USER
ARG MYSQL_PASSWORD
ARG ECHO_PORT
ARG ALLOW_ORIGIN

COPY ./ /var/www/work-counter/backend

WORKDIR /var/www/work-counter/backend

RUN apk add --no-cache gcc libc-dev && go build && go get github.com/rubenv/sql-migrate/...

EXPOSE 1313

CMD go run main.go