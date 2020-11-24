
.PHONY: all

all: pre_commit

.PHONY: test
test:
	mkdir -p .test-result
	go test -cover -coverprofile cover.out -outputdir .test-result ./...
	go tool cover -html=.test-result/cover.out -o .test-result/coverage.html

.PHONY: pre_commit
pre_commit: test
	go fmt ./...
	go vet ./...
