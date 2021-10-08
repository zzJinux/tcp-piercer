export DOCKER := docker
export COMPOSE := docker-compose
export MODULE_PATH := $(shell go list -m)
export GOMODCACHE

TEST_MAIN := test/utils/test_main.sh

.PHONY: test unittest composetest

test: unittest composetest

composetest:
	$(TEST_MAIN) $(TEST_OPTS) proxy_net_test

unittest:
	go test ./share/message
