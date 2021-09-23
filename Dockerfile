FROM golang:1.17-alpine as builder

WORKDIR /go/src/greenlight

COPY go.mod go.sum ./
RUN go mod download
RUN go get gotest.tools/gotestsum

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/api -x /go/src/greenlight/cmd/api

EXPOSE 4000

CMD ["/go/bin/api"]