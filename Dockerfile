FROM golang:1.13
WORKDIR /work
ADD pkg .
ADD server.go .
ADD go.mod .
RUN go build
CMD ["./image-server"]
