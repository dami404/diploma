PATH_TO_PROTO_DIR=proto
PATH_TO_PROTO_FILE=proto/auth/auth.proto
OUT=./gen
OPT=source_relative
PROTOC_PACKAGE=google.golang.org/protobuf/cmd/protoc-gen-go@latest
PROTOC_GEN_PACKAGE=google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
ENV_PATH="$PATH:$(go env GOPATH)/bin"

all:
	go install $(PROTOC_PACKAGE) && \
	go install $(PROTOC_GEN_PACKAGE) && \
	export PATH=$(PATH) &&  \
	protoc -I $(PATH_TO_PROTO_DIR) $(PATH_TO_PROTO_FILE) --go_out=$(OUT) --go_opt=paths=$(OPT) --go-grpc_out=$(OUT) --go-grpc_opt=paths=$(OPT)