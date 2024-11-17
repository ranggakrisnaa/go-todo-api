FROM golang:1.23.2 AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

COPY .env ./.env

WORKDIR /build/app

RUN CGO_ENABLED=0 go build -o app main.go

FROM alpine:latest

RUN apk update && apk add --no-cache tzdata

WORKDIR /app

COPY --from=build /build/app/app .

COPY --from=build /build/.env /app/.env

RUN chmod +x /app/app

EXPOSE 9090

CMD ["/app/app"]
