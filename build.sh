#!/bin/sh -xv

export GOPATH=`pwd`

go get -u github.com/mattn/go-sqlite3
go get -u github.com/nathanwinther/go-awsses
go get -u github.com/nathanwinther/go-uuid4

go install timesheet
go install scratch

