FROM golang:1.24-alpine AS builder
WORKDIR /go/src/github.com/dasLampe/recipientVerification.rwth/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

# Small image to run programm
FROM scratch
COPY --from=builder /app /app
CMD ["/app"]
