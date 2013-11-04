#!/bin/bash

if [ -z "$GOPATH" ]; then
    echo "Need to set GOPATH"
    exit 1
fi  
 
goget(){
	printf "Installing $1 : "
	go get $1

	if [ $? -eq 0 ]; then
	    echo OK
	else
	    echo FAIL
	    exit 1
	fi
}

goget github.com/PuerkitoBio/gocrawl


