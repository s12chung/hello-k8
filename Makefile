run:
	go install
	$(GOPATH)/bin/hello-k8

lint:
	golangci-lint run ./main.go
	golangci-lint run ./go/...

test:
	go test ./go/...

test-report:
	go test -v -covermode=atomic -coverprofile=coverage.out ./go/...
