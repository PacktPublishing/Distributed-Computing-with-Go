FROM golang:1.10


ADD . /go/src/littlefs

WORKDIR /go/src/littlefs

RUN go install littlefs

ENTRYPOINT /go/bin/littlefs

