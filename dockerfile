FROM golang

WORKDIR $GOPATH/src/github.com/asalih/bin-magicbytes

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o .

EXPOSE 6060

CMD ["bin-magicbytes"]