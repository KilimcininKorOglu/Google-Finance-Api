FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY cmd/ cmd/
COPY internal/ internal/
COPY web/ web/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/google-finance-api ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/google-finance-api /usr/local/bin/google-finance-api

EXPOSE 8080

ENTRYPOINT ["google-finance-api"]
