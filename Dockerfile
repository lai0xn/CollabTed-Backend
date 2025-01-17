FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go run github.com/steebchen/prisma-client-go generate --schema ./prisma
RUN go run github.com/steebchen/prisma-client-go db push --schema ./prisma

EXPOSE 8080

CMD ["go","run","cmd/server/main.go"]