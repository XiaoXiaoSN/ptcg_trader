MAKEFLAGS += -j5

PRJ_PATH = $(PWD)
GOTEST = $(go test -v)

.PHONY: server matcher build.docker swagger.server swagger.gen


##############################
# run service
##############################

server matcher:
	PROJ_HOME=$(CURDIR) go run main.go $@

docker.build:
	docker build . -f deploy/docker/trader.dockerfile -t xiao4011/ptcg_trader

docker.push:
	docker push xiao4011/ptcg_trader

docker.server:
	docker-compose up -d --build --force-recreate \
		--scale trader=3 \
		nginx trader matcher postgres redis swagger stan

gencode: swagger.gen mock.gen


##############################
# lint code
##############################

.PHONY: lint.prepare lint.go lint.gosec lint.all

lint.prepare:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

lint.go:
	golangci-lint run ./...

lint.sec:
	gosec ./...

lint.all: lint.go lint.sec


##############################
# swagger api document
##############################

swagger.server:
	docker run -d --rm --name ptcg_swagger -p 8088:8080 -e SWAGGER_JSON=/documents/swagger.json -v $(PRJ_PATH)/documents:/documents swaggerapi/swagger-ui

SWAGGER_FILE := documents/swagger.json
API_HEADER_FILE := $(PRJ_PATH)/pkg/delivery/restful/router.go
API_PATH := $(PRJ_PATH)/pkg
swagger.gen:
	# go get -u github.com/mikunalpha/goas
	goas --module-path . --main-file-path $(API_HEADER_FILE) --handler-path $(API_PATH) --output $(SWAGGER_FILE)


##############################
# mocking test data
##############################

mock.gen: mock.gen.svc.matcher mock.gen.repo mock.gen.redis

mock.gen.svc.matcher:
	# go install github.com/vektra/mockery/v2@latest
	mockery --dir pkg/service --name Matcher --filename mock_svc_matcher.go --structname MockMatcher --output test/mocks

mock.gen.repo:
	# go install github.com/vektra/mockery/v2@latest
	mockery --dir pkg/repository --name Repositorier --filename mock_repository.go --structname MockRepository --output test/mocks

mock.gen.redis:
	# go install github.com/vektra/mockery/v2@latest
	mockery --dir internal/redis --name Redis --filename mock_redis.go --structname MockRedis --output test/mocks
