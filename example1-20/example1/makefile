goVersion = $(shell go version)
gitHash = $(shell git show -s --format=%H)
buildTime = $(shell date "+%Y%m%d %T")

all:
	go build -ldflags "-X 'main.goVersion=$(goVersion)' -X 'main.gitHash=$(gitHash)' -X 'main.buildTime=$(buildTime)'" -o main

clean:
	rm -rf main
