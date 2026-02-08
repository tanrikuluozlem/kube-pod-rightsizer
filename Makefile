BINARY=kube-pod-rightsizer
IMAGE=kube-pod-rightsizer
TAG?=latest

.PHONY: build test lint docker clean

build:
	go build -o $(BINARY) ./cmd

test:
	go test -v ./...

lint:
	golangci-lint run

docker:
	docker build -t $(IMAGE):$(TAG) .

clean:
	rm -f $(BINARY)

run:
	go run ./cmd

deploy:
	kubectl apply -f deploy/

undeploy:
	kubectl delete -f deploy/
