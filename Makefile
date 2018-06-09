PACKAGES ?= $(shell go list ./... | grep -v /vendor/ | grep -v /tests)

.PHONY: all
all: vet lint megacheck test

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: lint
lint:
	STATUS=0; for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || STATUS=1; done; exit $$STATUS

.PHONY: megacheck
megacheck:
	STATUS=0; for PKG in $(PACKAGES); do CGO_ENABLED=0 megacheck $$PKG || STATUS=1; done; exit $$STATUS

.PHONY: test
test:
	STATUS=0; for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || STATUS=1; done; exit $$STATUS

.PHONY: protoc
protoc: noun/noun.proto gender/gender.proto case/case.proto number/number.proto
	protoc -I . gender/gender.proto --go_out=plugins=grpc:${GOPATH}/src
	protoc -I . case/case.proto --go_out=plugins=grpc:${GOPATH}/src
	protoc -I . number/number.proto --go_out=plugins=grpc:${GOPATH}/src
	protoc -I . noun/noun.proto --go_out=plugins=grpc:${GOPATH}/src
