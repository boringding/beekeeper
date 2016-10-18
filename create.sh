#!/bin/sh

if [ ! $1 ]; then
	echo "Usage: ./create.sh <project name>"
	exit
fi

projName=$1

cp -r ./skeleton $GOPATH/$projName

cd $GOPATH/$projName/conf

sed -i "s/beekeeper/$projName/g" framework.conf.dev.xml
sed -i "s/beekeeper/$projName/g" framework.conf.test.xml
sed -i "s/beekeeper/$projName/g" framework.conf.product.xml

cd $GOPATH/$projName

sed -i "s/github.com\/boringding\/beekeeper\/skeleton\/handlers/$projName\/handlers/g" main.go