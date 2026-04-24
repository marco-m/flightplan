#! /bin/sh

set -ex
pwd
ls
VERSION=$(date +%Y%m%d%H%M%S)
GIFT=gift/gift-$VERSION
echo "hello" > $GIFT
