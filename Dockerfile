FROM golang:1.23

WORKDIR /app

COPY . .
RUN go build -o 1337b04rd cmd/main.go

EXPOSE 8080

CMD ["./1337b04rd"]
