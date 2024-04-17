gen:
	@protoc \
		--proto_path=protobuf "protobuf/nodes.proto" \
		--go_out=services/genproto/nodes --go_opt=paths=source_relative \
  	--go-grpc_out=services/genproto/nodes --go-grpc_opt=paths=source_relative