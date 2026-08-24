FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go mod tidy && go build -o gateway ./cmd/gateway

FROM alpine:3.20
COPY --from=builder /app/gateway /usr/local/bin/
CMD ["gateway"]
