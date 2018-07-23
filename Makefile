PACKAGES ?= $(shell go list ./... | grep -v /vendor/ | grep -v /tests)

.PHONY: all
all: vet lint megacheck test

vet:
	go vet $(PACKAGES)

lint:
	STATUS=0; for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || STATUS=1; done; exit $$STATUS

megacheck:
	STATUS=0; for PKG in $(PACKAGES); do CGO_ENABLED=0 megacheck $$PKG || STATUS=1; done; exit $$STATUS

test:
	STATUS=0; for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || STATUS=1; done; exit $$STATUS

generate:
	go generate ./...


build-%:
	docker build -f Dockerfile.$* -t binquiry/$* .
	docker tag binquiry/$* localhost:5000/binquiry/$*

push-%:
	docker push localhost:5000/binquiry/$*

run-nounparser:
	docker run --rm --name nounparser -d -p 50051:50051 binquiry/nounparser

run-frontend:
	docker run --rm --link nounparser:nounparser --name frontend -d -p 8000:8000 -e NOUNPARSER=nounparser:50051 -e LOGLEVEL=debug binquiry/frontend

stop:
	docker container stop frontend nounparser

k8s-export-%:
	kubectl get --export -o json deployments/$* > k8s/deployment-$*.yaml
	kubectl get --export -o json services/$* > k8s/service-$*.yaml

k8s-create:
	kubectl create -f k8s/deployment-nounparser.yaml
	kubectl create -f k8s/service-nounparser.yaml
	kubectl create -f k8s/deployment-frontend.yaml
	kubectl create -f k8s/service-frontend.yaml

k8s-delete:
	kubectl delete -f k8s/service-frontend.yaml
	kubectl delete -f k8s/deployment-frontend.yaml
	kubectl delete -f k8s/service-nounparser.yaml
	kubectl delete -f k8s/deployment-nounparser.yaml

k8s-test-%:
	curl "http://`minikube ip`:`kubectl get services/frontend -o go-template='{{(index .spec.ports 0).nodePort}}'`/$*"
