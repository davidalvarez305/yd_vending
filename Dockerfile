# syntax=docker/dockerfile:1

FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV APP_ENV=production

RUN go build -o /exec

CMD [ "/exec" ]