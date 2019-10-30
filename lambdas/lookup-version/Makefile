all: lookup-version

lookup-version: main.go
	go build

deployable: main.go
	GOOS=linux GOARCH=amd64 go build
	chmod a+x lookup-version
	zip ../deployables/lookup-version.zip lookup-version
	rm -f lookup-version

test: main.go
	go test

clean:
	rm -f lookup-version
