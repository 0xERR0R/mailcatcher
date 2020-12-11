FROM golang:1.15-alpine AS builder

RUN apk --no-cache add gcc

WORKDIR ${GOPATH}/src/github.com/0xERR0R/mailcatcher
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/mailcatcher cmd/mailcatcher/*.go

FROM alpine

LABEL org.opencontainers.image.source="https://github.com/0xERR0R/mailcatcher" \
      org.opencontainers.image.url="https://github.com/0xERR0R/mailcatcher" \
      org.opencontainers.image.title="Self hosted mail trash service"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/mailcatcher /app
ENTRYPOINT ["/app"]
