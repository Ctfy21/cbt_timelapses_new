FROM golang:latest

COPY .. .

RUN go mod download
RUN go build -v -o main ./cmd/app/

CMD ["./main"]