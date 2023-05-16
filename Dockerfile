########-- Build stage --########

FROM golang:1.20-alpine AS builder-prod

RUN  apk add git

WORKDIR /opt

COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY . ./
RUN go build -o /atackhelp ./cmd


########-- Deploy stage --########

FROM alpine:3.18

WORKDIR /opt 

COPY --from=builder-prod /atackhelp /opt/atackhelp
COPY ./.env .

ARG APP_PORT

CMD [ "/opt/atackhelp"]

EXPOSE $APP_PORT
