FROM golang:alpine

RUN apk update && apk add --no-cache 

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["air", "-c", ".air.toml"]
