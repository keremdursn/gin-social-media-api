FROM golang:1.21-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
