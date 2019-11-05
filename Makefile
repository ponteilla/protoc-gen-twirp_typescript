TESTABLE=$$(go list ./...)

all: test install

deps:
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure
	@go get -u github.com/twitchtv/retool
	@retool sync
.PHONY: deps

install:
	@go install
.PHONY: install

gen: install
	@retool do go generate ./...
.PHONY: gen

test:
	@go test -v $(TESTABLE)
.PHONY: test
