#!/usr/bin/env bash

set -ex
echo "build gochat.bin ..."
# CGO_CFLAGS="-g -O2 -Wno-return-local-addr" fix compile sqlite3 warning
CGO_CFLAGS="-g -O2 -Wno-return-local-addr" go build -o /go/src/gochat/bin/gochat.bin /go/src/gochat/main.go
echo "restart all ..."
supervisorctl restart all
echo "all Done."
echo "Beautiful ! Now, You can visit http://127.0.0.1:8080 , start the world."