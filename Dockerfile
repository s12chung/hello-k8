FROM golang:1.10.4-alpine3.8

RUN apk --no-cache add\
    make git dep

ENV PORT 8080

WORKDIR /go/src/github.com/s12chung/hello-k8
COPY . .

CMD ["tail", "-f", "/dev/null"]
