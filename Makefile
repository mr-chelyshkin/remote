.PHONY: build_darwin build_linux build_arm

build_darwin:
	GOOS=darwin go build -o remote ./cmd

build_linux:
	GOOS=linux go build -o remote ./cmd

build_arm:
	GOOS=linux GOARM=7 GOARCH=arm go build -o remote ./cmd
