#!/bin/bash

# Your PNGs will be here
GIFDIR=/tmp/pepparkak

SCRIPTDIR=$(cd `dirname $0` && pwd)
ROOTDIR=$(cd "${SCRIPTDIR}/.." && pwd)

cd "${ROOTDIR}"

mkdir -p "${GIFDIR}"

go build cmd/pepparkak/pepparkak.go
./pepparkak ${GIFDIR}

