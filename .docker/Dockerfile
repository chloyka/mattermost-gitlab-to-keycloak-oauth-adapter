FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN --mount=target=/go,type=cache go mod download
RUN --mount=target=/go,type=cache CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o binary .

RUN chmod +x binary

FROM scratch as run
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/binary /binary
