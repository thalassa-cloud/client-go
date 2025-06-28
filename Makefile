PKG_LIST := $(shell go list ./... | grep -v tests)

default: test

test: 
	@go test -short ${PKG_LIST}            

review:
	reviewdog -diff="git diff FETCH_HEAD" -tee
