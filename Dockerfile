FROM golang

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

EXPOSE 3001

RUN /go/bin/app