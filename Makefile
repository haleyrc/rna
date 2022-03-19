all: build run

.PHONY: clean
clean:
	@rm -f build/rna
	@rm -rf test

.PHONY: build
build:
	@go build -o build/rna .

.PHONY: run
run:
	@./build/rna

.PHONY: test
test: clean build
	@./build/rna new test --name test
