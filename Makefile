.PHONY: generate
generate:
	go fmt
	go vet
	#go generate

.PHONY: build
build: generate
	go build -o bin/vault-agent
