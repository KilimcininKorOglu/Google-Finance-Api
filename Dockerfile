FROM golang:1.22-alpine@sha256:1699c10032ca2582ec89a24a1312d986a3f094aed3d5c1147b19880afe40e052 AS builder

WORKDIR /src
COPY go.mod ./
COPY cmd/ cmd/
COPY internal/ internal/
COPY web/ web/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/google-finance-api ./cmd/server

FROM alpine:3.20@sha256:d9e853e87e55526f6b2917df91a2115c36dd7c696a35be12163d44e6e2a4b6bc

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S -G app app

COPY --from=builder /bin/google-finance-api /usr/local/bin/google-finance-api

USER app

EXPOSE 8080

ENTRYPOINT ["google-finance-api"]
