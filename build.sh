#!/bin/bash

function build_once()
{
    export GOOS=$1
    export GOARCH=$2
	export TAG=$3

    echo build $GOOS $GOARCH

    export DIR=autoproxy_"$GOOS"_"$GOARCH"
    mkdir $DIR

	cp -rf domain.json $DIR/domain.json
	go build -ldflags="-w -s" -o $DIR/autoproxy$TAG .
    zip -r -q -o autoproxy_"$GOOS"_"$GOARCH".zip $DIR/*
	rm -rf $DIR
}

rm -rf ./*.zip

build_once windows amd64 .exe
build_once windows 386 .exe
build_once darwin amd64
build_once darwin arm64
build_once linux amd64
build_once linux 386
build_once linux arm64
build_once linux arm
