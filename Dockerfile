########-- Build stage --########

FROM golang:1.20-alpine AS builder-prod

WORKDIR /opt

COPY . /opt

RUN go build -mod=vendor -o ./runner ./cmd

########-- Deploy stage --########

FROM alpine:3.18

WORKDIR /opt 

COPY --from=builder-prod /opt/runner /opt/.env /opt/
COPY ./app/services/template /opt/app/services/template

ENTRYPOINT /opt/runner