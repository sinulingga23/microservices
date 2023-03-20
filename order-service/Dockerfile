FROM golang:1.19-alpine AS builder
RUN apk add --no-cache ca-certificates git

WORKDIR /opt/apps

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /order-service .

FROM alpine as release
RUN apk add --no-cache ca-certificates

WORKDIR /order-service
COPY --from=builder /order-service ./main
ENTRYPOINT ["/order-service/main"]