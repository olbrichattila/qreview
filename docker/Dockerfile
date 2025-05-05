FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app


RUN go install github.com/olbrichattila/godbmigrator_cmd/cmd/migrator@latest


FROM alpine:latest


COPY --from=builder /go/bin/migrator /usr/local/bin/migrator


ENTRYPOINT ["/usr/local/bin/migrator"]
