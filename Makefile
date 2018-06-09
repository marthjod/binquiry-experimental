protoc: noun/noun.proto gender/gender.proto case/case.proto number/number.proto
	protoc -I . gender/gender.proto --go_out=plugins=grpc:${GOPATH}/src
	protoc -I . case/case.proto --go_out=plugins=grpc:${GOPATH}/src
	protoc -I . number/number.proto --go_out=plugins=grpc:${GOPATH}/src
	protoc -I . noun/noun.proto --go_out=plugins=grpc:${GOPATH}/src
