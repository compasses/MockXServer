#!/bin/bash
version=`git log --date=iso --pretty=format:"%cd @%h" -1`
if [ $? -ne 0 ]; then
    version="not a git repo"
fi

compile=`date +"%F %T %z"`" by "`go version`

cat << EOF | gofmt > version.go
package main

const (
    Version = "$version"
    Compile = "$compile"
)
EOF
