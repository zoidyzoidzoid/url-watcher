generate:
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=plugins=grpc:pkg/ \
	  --grpc-gateway_out=logtostderr=true:pkg/ \
	  --swagger_out=logtostderr=true:docs/ \
	  proto/watcher_service.proto
