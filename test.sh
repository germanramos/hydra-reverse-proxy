#!/bin/sh -e

. ./build

go test -i ./reverse_proxy
go test -v ./reverse_proxy

go test -i ./tests/acceptance
go test -v ./tests/acceptance