#!/bin/bash
currentpath=`pwd`
dest="/usr/local/pipe"

mkdir -p ${dest}
mkdir -p ${dest}/{conf,bin,run,www}

cp -f ${currentpath}/bin/pipe-server -p ${dest}/bin/
cp -f ${currentpath}/src/config/config.yaml -p ${dest}/conf/
