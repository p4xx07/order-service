FROM golang:1.23-alpine AS build_base

WORKDIR /go/src/app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) go build -ldflags="-s -w"  -o order-service ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build_base /go/src/app/order-service /app/order-service

EXPOSE 8080
CMD ["./order-service"]
