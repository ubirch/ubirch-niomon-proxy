######################################################################
# @author      : annika (annika@berlin.ccc.de)
# @file        : Makefile
# @created     : Sunday Aug 16, 2020 19:24:54 CEST
######################################################################

all: ubproxy

ubproxy: build

build: 
	cd cmd/ubproxy && go build
	mkdir -p ./bin
	mv cmd/ubproxy/ubproxy ./bin
	chmod 755 bin/ubproxy

.PHONY:
	clean test

run:
	go run cmd/ubproxy/main.go

runbin:
	./bin/ubproxy

test:
	cd pkg/ubproxy/model/expkey && go test -v
	cd pkg/ubproxy/expkeyservice && go test -v

gitupdate:
	git pull

clean:
	rm -f bin/*
	rm -f cmd/ubproxy/docs/*


refreshrun: clean gitupdate build runbin
