FROM golang:latest AS builder
WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV CGO_ENABLED 0
RUN go build -o /assets/check ./cmd/check \
    && go build -o /assets/in ./cmd/in \
    && go build -o /assets/out ./cmd/out

RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM ubuntu:bionic AS resource
ENV DEBAIN_FRONTEND=noninteractive
RUN apt-get update
COPY --from=builder assets/ opt/resource/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN chmod +x /opt/resource/*

FROM resource AS tests
COPY --from=builder /tests /tests

ARG WEBHOOK_ID
ARG WEBHOOK_TOKEN

RUN set -e; for test in /tests/*.test; do \
		$test -ginkgo.v; \
	done

FROM resource
