FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/main
RUN CGO_ENABLED=0 GOOS=linux go build -o main
WORKDIR /app
EXPOSE 8085
CMD ["./cmd/main/main"]
