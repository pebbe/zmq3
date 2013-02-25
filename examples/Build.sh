#!/bin/sh

dirs=''
for i in *
do
    if [ -d $i ]
    then
	dirs="$dirs $i"
    fi
done

libs=`pkg-config --libs-only-L libzmq`
if [ "$libs" = "" ]
then
    for i in $dirs
    do
	cd $i
	go build .
	cd ..
    done
else
    libs="-r `echo $libs | perl -p -e 's/-L//; s/ -L/:/g'`"
    for i in $dirs
    do
	cd $i
	go build -ldflags="$libs" .
	cd ..
    done
fi
