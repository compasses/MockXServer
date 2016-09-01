#!/bin/bash

set -ex

# Create go binary and package verifier + mock service into distribution
VERSION=$(go version)
echo "==> Go version ${VERSION}"

echo "==> Getting dependencies..."
export GO15VENDOREXPERIMENT=1

go get github.com/mitchellh/gox
go get github.com/inconshreveable/mousetrap # windows dep
go get -d ./...
gox -os="darwin" -arch="amd64" -output="cmd/MockXServer_{{.OS}}_{{.Arch}}"
gox -os="windows" -arch="386" -output="cmd/MockXServer_{{.OS}}_{{.Arch}}"
gox -os="linux" -arch="386" -output="cmd/MockXServer_{{.OS}}_{{.Arch}}"
gox -os="linux" -arch="amd64" -output="cmd/MockXServer_{{.OS}}_{{.Arch}}"

echo
echo "==> Results:"
ls -hl cmd/
#echo "packaging ..."
#tar -czf MockXServer.tar.gz MockXServer/
