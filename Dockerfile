ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine

RUN apk add bash curl

WORKDIR /source

COPY offergen/go.mod offergen/go.sum /source/
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@v0.3.833
RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.22.2
RUN go install github.com/air-verse/air@v1.61.7

ENTRYPOINT [ "go" ]
