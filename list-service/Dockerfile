FROM golang:alpine as build-env

RUN apk add git
RUN mkdir /simpletodo
RUN apk --update add ca-certificates
WORKDIR /simpletodo
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/simpletodo-list-service
ENTRYPOINT ["/go/bin/simpletodo-list-service"]