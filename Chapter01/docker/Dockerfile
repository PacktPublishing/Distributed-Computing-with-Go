FROM golang:1.9.1
# The base image we want to use to build our docker image from.
# Since this image is specialized for golang it will have GOPATH = /go


ADD . /go/src/hello
# We copy files & folders from our system onto the docker image

RUN go install hello
# Next we can create an executable binary for our project with the command, `go install`

ENV NAME Bob

ENTRYPOINT /go/bin/hello
# Command to execute when we start the container

