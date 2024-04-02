build:
	go build -ldflags="-s -w" goctl.go
	$(if $(shell command -v upx), upx goctl)

mac:
	GOOS=darwin go build -ldflags="-s -w" -o goctl-darwin goctl.go
	$(if $(shell command -v upx), upx goctl-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o goctl.exe goctl.go
	$(if $(shell command -v upx), upx goctl.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o goctl-linux goctl.go
	$(if $(shell command -v upx), upx goctl-linux)

image:
	docker build --rm --platform linux/amd64 -t kevinwan/goctl:$(version) .
	docker tag kevinwan/goctl:$(version) kevinwan/goctl:latest
	docker push kevinwan/goctl:$(version)
	docker push kevinwan/goctl:latest
	docker build --rm --platform linux/arm64 -t kevinwan/goctl:$(version)-arm64 .
	docker tag kevinwan/goctl:$(version)-arm64 kevinwan/goctl:latest-arm64
	docker push kevinwan/goctl:$(version)-arm64
	docker push kevinwan/goctl:latest-arm64



rpc-gen:
	./goctl rpc protoc testout/*.proto --go_out=testout/rpc --go-grpc_out=testout/rpc --zrpc_out=testout/rpc -m -home=/Users/xiaohan/soft/project/xiaohan/infra/templates -style=go_zero -w=true

model-gen:
	./goctl model mysql ddl --src *.sql -c --dir=./internal/model -home=/Users/x/soft/project/y/infra/templates -style=go_zero -prefix=infra -w=true

api-gen:
	./goctl api go --api=testout/test.api --dir=testout/api -home=/Users/xiaohan/soft/project/xiaohan/infra/templates -style=go_zero -w=true