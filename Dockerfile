# Image to build program.
FROM golang:1.14 AS builder
WORKDIR /go/src/git.alania-breslau.de/server/empfaengerverifizierung/
COPY empfaengerverifizierung.go .

RUN go get -d -v ./
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Small image to run programm
FROM alpine:latest
WORKDIR /opt

COPY templates/ ./templates/
COPY .env .

COPY --from=builder /go/src/git.alania-breslau.de/server/empfaengerverifizierung/app .

CMD ["./app"]