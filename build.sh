#!/bin/sh
set -ex
APP=mos
VERSION=${VERSION:-0.3}
APPDIR=${DESTDIR:-${APP}_$VERSION}
CGO_ENABLED=0
DATE="$(date -u '+%Y-%m-%d %I:%M:%S %Z')"
COMMIT=$(git rev-parse HEAD)

mkdir -p $APPDIR/usr/bin

go build -ldflags "-X \"github.com/MocaccinoOS/mos-cli/cmd.Version=$VERSION\" \
          -X \"github.com/MocaccinoOS/mos-cli/cmd.BuildTime=$DATE\" \
          -X \"github.com/MocaccinoOS/mos-cli/cmd.BuildCommit=$COMMIT\"" \
         -o $APPDIR/usr/bin/$APP main.go