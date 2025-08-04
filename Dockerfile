FROM golang:1.23.3-alpine

WORKDIR /app/cmd

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

CMD ["./main"]