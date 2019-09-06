FROM golang:1.13
WORKDIR /work
ADD . .
RUN go build
CMD ["./image-server"]
