FROM golang:1.10


ADD . /go/src/github.com/last-ent/distributed-go/chapter8/goophr/librarian

RUN go install github.com/last-ent/distributed-go/chapter8/goophr/librarian

ENTRYPOINT /go/bin/librarian

