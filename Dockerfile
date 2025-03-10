# syntax=docker/dockerfile:1

FROM golang:1.22.5 AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /output/ppparsser

# --
FROM alpine
RUN apk add curl
WORKDIR /app
COPY --from=builder /output/ppparsser ./ppparsser
COPY data/* ./data/

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD curl --fail http://localhost:8080 || exit 1
CMD ["/app/ppparsser"]
