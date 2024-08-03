# syntax=docker/dockerfile:1

FROM golang:1.22.5 AS builder
WORKDIR /app

COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /ppparsser

# --
FROM alpine
RUN apk add curl
WORKDIR /app
COPY --from=builder /test-result /ppparsser

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD curl --fail http://localhost:8080 || exit 1
CMD ["/ppparsser"]
