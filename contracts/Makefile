PATH_TO_PROTO_DIR=proto
PATH_TO_PROTO_FILE=proto/auth/auth.proto
OUT=./gen
OPT=source_relative



all:
	protoc -I $(PATH_TO_PROTO_DIR) $(PATH_TO_PROTO_FILE) --go_out=$(OUT) --go_opt=paths=$(OPT) --go-grpc_out=$(OUT) --go-grpc_opt=paths=$(OPT)