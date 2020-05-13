FROM golang:1.13-alpine

COPY ./backend /var/www/work-counter/backend

WORKDIR /var/www/work-counter/backend

RUN apk add --no-cache gcc libc-dev && go build && go get github.com/rubenv/sql-migrate/...

EXPOSE 1313

CMD go run main.go