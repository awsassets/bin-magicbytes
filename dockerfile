FROM golang

RUN go get -x github.com/asalih/bin-magicbytes
WORKDIR /go/src/github.com/asalih/bin-magicbytes
RUN make


EXPOSE 6060