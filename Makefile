MAKEFLAGS += -j5

PRJ_PATH = $(PWD)
GOTEST = $(go test -v)

server:
	PROJ_HOME=$(CURDIR) CONFIG_NAME=app go run main.go $@

build.docker:
	docker build . -f deploy/docker/trader.dockerfile
