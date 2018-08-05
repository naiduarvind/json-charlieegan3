.PHONY: test

test:
	go test $$(go list ./...)

image:
	docker build -t charlieegan3/json-charlieegan3:$$(tar -cf - . | md5sum | awk '{ print $$1 }') .
	docker push charlieegan3/json-charlieegan3:$$(tar -cf - . | md5sum | awk '{ print $$1 }')
