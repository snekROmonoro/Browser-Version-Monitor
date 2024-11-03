FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go run github.com/steebchen/prisma-client-go generate

CMD ["go", "run", "main.go"]
