FROM golang:1.10.4-alpine3.8

# first line for basic dev
# second line for goose migrations
RUN apk --no-cache add\
    make git dep\
    gcc libc-dev

ENV PORT 8080

RUN go get -u github.com/pressly/goose/cmd/goose

WORKDIR /go/src/github.com/s12chung/hello-k8
COPY . .

CMD ["tail", "-f", "/dev/null"]
