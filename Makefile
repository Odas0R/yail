build:
	go build -o yail .

run-ast:
	go build -o yail . && ./yail ast

run-vm:
	go build -o yail . && ./yail vm

:PHONY: build repl-ast repl-vm
