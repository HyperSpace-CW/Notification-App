FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main ./cmd/app/main.go

CMD ["./main"]

EXPOSE 3003