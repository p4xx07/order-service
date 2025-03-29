FROM golang:1.23-alpine AS build_base

WORKDIR /go/src/app

COPY ./go.mod .
COPY ./go.sum .
COPY ./.env .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"  -o order-service ./main.go

FROM phusion/baseimage:focal-1.2.0

WORKDIR /app

COPY --from=build_base /go/src/app/order-service /app/order-service

EXPOSE 8080
CMD ["./order-service"]
