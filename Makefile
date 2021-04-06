

default: build_linux image clean

clean:
	rm -f mdout

build_mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build .

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

image:
	docker build -t mdout .
	docker build -t mdout:chinese --build-arg LANGUAGE=chinese .