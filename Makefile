all: build run

.PHONY: build
build:
	@go build -o build/rna .

.PHONY: run
run:
	@./build/rna
