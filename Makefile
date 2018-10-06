build:
	docker build -t hello-k8:v1.0 .

apply:
	minikube mount .:/dev-mount &>/dev/null &disown;
	kubectl apply -f ./manifests/
	@echo
	@echo hello-k8 is being served at:
	@echo $$(minikube service hello-k8 --url)
	@echo

# kubectl describe ...
status:
	kubectl get all
	kubectl get pv
	kubectl get pvc

clean:
	pkill -f "minikube mount" || true
	kubectl delete deployments.app --all
	kubectl delete service --all

clean-all: clean
	kubectl delete pvc --all
	kubectl delete pv --all
	kubectl delete secrets --all

install:
	dep ensure

exec-sh:
	kubectl exec $(shell kubectl get pods -o=custom-columns=name:metadata.name | grep hello-k8) -ti ash

run:
	go install
	$(GOPATH)/bin/hello-k8

lint:
	golangci-lint run ./main.go
	golangci-lint run ./go/...

test:
	go test ./go/...

test-apply:
	kubectl apply -f ./manifests/hello-k8.test.yaml
	kubectl apply -f ./manifests/postgres.yaml

# Using cat with coverage.out instead of kubectl cp because kubectl cp can't find the file for some reason in Travis CI
test-ci:
	kubectl exec $(shell kubectl get pods -o=custom-columns=name:metadata.name | grep hello-k8) -ti make db-up
	kubectl exec $(shell kubectl get pods -o=custom-columns=name:metadata.name | grep hello-k8) -ti make test-report
	kubectl exec $(shell kubectl get pods -o=custom-columns=name:metadata.name | grep hello-k8) -ti cat coverage.out > coverage.out

test-report: install
	go test -v -covermode=atomic -coverprofile=coverage.out ./go/...

DB_STRING := "user=postgres password=$(POSTGRES_PASSWORD) dbname=postgres host=$(POSTGRES_SERVICE_HOST) port=$(POSTGRES_SERVICE_PORT) sslmode=disable"

db-up:
	goose -dir ./migrations postgres $(DB_STRING) up

db-down:
	goose -dir ./migrations postgres $(DB_STRING) down

db-redo:
	goose -dir ./migrations postgres $(DB_STRING) redo

db-status:
	goose -dir ./migrations postgres $(DB_STRING) status
