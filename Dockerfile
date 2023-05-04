FROM golang

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...
RUN /go/bin/app