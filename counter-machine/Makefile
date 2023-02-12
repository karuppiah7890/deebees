# TODO: Build using https://goreleaser.com/
build:
	CGO_DISABLE=0 go build -v

build-linux:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -v

release-linux: build-linux
	tar cvzf counter-machine-linux-amd64.tar.gz counter-machine

# TODO: Lint using golangci-lint
