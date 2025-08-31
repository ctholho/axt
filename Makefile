ARGS ?=

.PHONY: test
test:
	TZ=UTC go test ./... $(ARGS)
