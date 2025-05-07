# Builder container
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app


RUN go install github.com/olbrichattila/qreview@latest

# Built container
FROM alpine:latest

COPY --from=builder /go/bin/qreview /usr/local/bin/qreview
COPY entrypoint.sh /entrypoint.sh
COPY definitions.yaml /definitions.yaml
RUN chmod +x /entrypoint.sh


ENTRYPOINT ["/entrypoint.sh"]
