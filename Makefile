export DOCKER := docker
export COMPOSE := docker-compose
export MODULE_PATH := $(shell go list -m)

TEST_MAIN := test/utils/test_main.sh

.PHONY: test

test:
	$(TEST_MAIN) $(TEST_OPTS) proxy_net_test
