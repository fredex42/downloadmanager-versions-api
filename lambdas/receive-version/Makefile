all: receive-version

receive-version: main.go
	go build

deployable: main.go
	GOOS=linux GOARCH=amd64 go build
	chmod a+x receive-version
	zip ../deployables/receive-version.zip receive-version
	rm -f receive-version

test: main.go
	go test

clean:
	rm -f receive-version
