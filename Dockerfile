ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine

RUN apk add bash curl

WORKDIR /source

COPY offergen/go.mod offergen/go.sum /source/
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@v0.2.747
RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.17.3
RUN go install github.com/cosmtrek/air@v1.52.0

ENTRYPOINT [ "go" ]
