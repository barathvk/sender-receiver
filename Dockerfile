FROM golang:1.16.2-alpine AS build
WORKDIR /app
COPY go.mod go.sum main.go /app/
COPY sender /app/sender
COPY receiver /app/receiver
COPY common /app/common
RUN go build .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/sender-receiver /app

ENTRYPOINT [ "./sender-receiver" ]