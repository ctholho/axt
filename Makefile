ARGS ?=

.PHONY: test
test:
	TZ=UTC go test ./... $(ARGS)

.PHONY: build
build:
	go build -o axt

.PHONY: debug
debug:
	go run debug/debug.go | ./axt $(ARGS)
