# Image to build program.
FROM golang:1.14 AS builder
WORKDIR /go/src/github.com/dasLampe/recipientVerification.rwth/
COPY recipientVerification.go .

RUN go get -d -v ./
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Small image to run programm
FROM alpine:latest
WORKDIR /opt

COPY templates/ ./templates/
COPY .env .

COPY --from=builder /go/src/github.com/dasLampe/recipientVerification.rwth/app .

CMD ["./app"]