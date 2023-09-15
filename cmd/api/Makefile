#init_api:
#	sh hz_gen.sh init

#update_api:
#	sh hz_gen.sh

.PHONY: init
init:
	hz new \
	-module github.com/SIN5t/tiktok_v2/cmd/api \
	-service "api-gateway" \
	-proto_path ../../ \
	-idl ../../idl/ApiGateway.proto
	&& go mod tidy

.PHONY: update
update:
	hz update \
	-proto_path ../../ \
	-idl ../../idl/ApiGateway.proto
.PHONY: run
run:
	go build -o hertz_test && ./hertz_test