all: yggd

.PHONY: yggd
yggd:
	go build -o bin/yggd ./cmd/yggd

.PHONY: up
up: yggd
	./bin/yggd
