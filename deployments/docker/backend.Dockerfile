FROM golang:1.26.0-alpine3.24 AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -o /app/bin/pulse-monolith ./cmd/monolith/main.go

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/pulse-monolith /pulse-monolith
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/pulse-monolith"]
