FROM golang:1.12
WORKDIR /go/src/github.com/kkty/image-server
ADD . .
CMD ["go", "run", "main.go"]
