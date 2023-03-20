FROM golang:1.19-alpine AS builder
RUN apk add --no-cache ca-certificates git

WORKDIR /opt/apps

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /product-service .

FROM alpine as release
RUN apk add --no-cache ca-certificates

WORKDIR /product-service
COPY --from=builder /product-service ./main
ENTRYPOINT ["/product-service/main"]