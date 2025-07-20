FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /build
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o netprobe ./cmd/netprobe

FROM scratch
COPY --from=builder /build/netprobe /netprobe
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER 65534:65534
ENTRYPOINT ["/netprobe"]
