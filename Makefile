.PHONY: test

test:
	go clean
	go build
	ginkgo
