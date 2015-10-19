GP := $(shell dirname $(realpath $(lastword $(GOPATH))))
ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export GOPATH := ${ROOT}/:${GOPATH}


atest:
	go test ./memorycache/


btest:
	go test ./src/

