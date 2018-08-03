.PHONY: test

test:
	go test $$(go list ./...)
