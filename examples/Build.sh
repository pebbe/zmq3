#!/bin/sh

libs=`pkg-config --libs-only-L libzmq`
if [ "$libs" = "" ]
then
    for i in *.go
    do
	go build $i
    done
else
    libs="-r `echo $libs | perl -p -e 's/-L//; s/ -L/:/g'`"
    for i in *.go
    do
	go build -ldflags="$libs" $i
    done

fi
