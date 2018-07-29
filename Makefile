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

run-dispatcher:
	docker run --rm --name dispatcher --link nounparser:nounparser -d -p 50052:50052 -e NOUNPARSER=nounparser:50051 binquiry/dispatcher

run-frontend:
	docker run --rm --name frontend --link dispatcher:dispatcher -d -p 8000:8000 -e DISPATCHER=dispatcher:50052 -e LOGLEVEL=debug binquiry/frontend

stop:
	docker container stop frontend nounparser

k8s-export-%:
	kubectl get --export -o json deployments/$* > k8s/deployment-$*.yaml
	kubectl get --export -o json services/$* > k8s/service-$*.yaml

k8s-create:
	for s in nounparser dispatcher frontend; do
		kubectl create -f k8s/deployment-${s}.yaml
		kubectl create -f k8s/service-${s}.yaml
	done

k8s-delete:
	for s in nounparser dispatcher frontend; do
		kubectl delete -f k8s/service-${s}.yaml
		kubectl delete -f k8s/deployment-${s}.yaml		
	done

k8s-test-%:
	@curl -s "http://`minikube ip`:`kubectl get services/frontend -o go-template='{{(index .spec.ports 0).nodePort}}'`/$*"
