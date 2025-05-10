FROM golang:1.24-alpine

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o myapp cmd/main.go

# RUN chmod +x myapp

CMD ["./myapp"]