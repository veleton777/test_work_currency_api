FROM golang:1.22-alpine3.20 AS builder

WORKDIR /work

COPY ../.. .

RUN apk add --update --no-cache git make bash
RUN make build-app

FROM alpine:3.20

RUN apk add --update --no-cache bash

RUN adduser -D appuser
USER appuser

WORKDIR /
COPY --from=builder /work/build/app .
COPY --from=builder /work/build/docker/wait-for-it.sh .
EXPOSE 8080
CMD ["/app", "http"]