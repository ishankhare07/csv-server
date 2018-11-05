FROM golang:1.11
RUN go get -u -v github.com/golang/dep/cmd/dep github.com/ishankhare07/csv-server
WORKDIR $GOPATH/src/github.com/ishankhare07/csv-server/
RUN dep ensure -v
CMD go run main.go
