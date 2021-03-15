FROM golang:1.16.2-alpine
WORKDIR /app
COPY go.mod go.sum main.go /app/
COPY sender /app/sender
COPY receiver /app/receiver
COPY common /app/common
RUN go build .
ENTRYPOINT [ "./sender-receiver" ]