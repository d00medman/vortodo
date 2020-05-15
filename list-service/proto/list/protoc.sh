protoc -I. \
	--include_imports \
    --include_source_info \
    --descriptor_set_out out.pb \
	--go_out=plugins=grpc:. \
	proto/list/list.proto