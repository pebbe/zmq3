#!/bin/sh

gobin=`go env GOBIN`
if [ "$gobin" = "" ]
then
    gobin=`go env GOPATH`
    if [ "$gobin" = "" ]
    then
	gobin=`go env GOROOT`
    fi
    gobin=`echo $gobin | sed -e 's/:.*//'`
    gobin=$gobin/bin
fi

dir=$gobin/zmq-examples

echo Installing examples in $dir

mkdir -p $dir

cp -u */*.sh $dir

dirs=''
for i in *
do
    if [ -d $i ]
    then
	if [ ! -f $dir/$i -o $dir/$i -ot $i/$i.go ]
	then
	    dirs="$dirs $i"
	fi
    fi
done

libs=`pkg-config --libs-only-L libzmq`
if [ "$libs" = "" ]
then
    for i in $dirs
    do
	cd $i
	go build -o $dir/$i .
	cd ..
    done
else
    libs="-r `echo $libs | sed -e 's/-L//; s/  *-L/:/g'`"
    for i in $dirs
    do
	cd $i
	go build -ldflags="$libs" -o $dir/$i .
	cd ..
    done
fi
