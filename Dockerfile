FROM golang:1.17.1

WORKDIR /go/src/linkservice
COPY . .
RUN go mod download
RUN go build -o ./cmd/linkservice/linkservice ./cmd/linkservice

EXPOSE 50051

CMD ["./cmd/linkservice/linkservice"]