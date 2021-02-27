FROM golang:1.15.6
WORKDIR /go/src/
RUN go mod download
RUN go build -o app calculator/server/*

FROM golang:1.15.6
WORKDIR /go/src/
RUN go build -o app calculator/client/*
CMD ["./app"]