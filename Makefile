build:
	go build -o yail .

run-ast:
	go build -o yail . && ./yail ast

run-vm:
	go build -o yail . && ./yail vm

test-cp:
	go build -o yail . && ./run-cp.sh

:PHONY: build repl-ast repl-vm test-cp
