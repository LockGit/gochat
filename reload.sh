#!/usr/bin/env bash

set -e

echo "replace api and websocket ip address..."
addr_http=${HOST_IP-127.0.0.1}:7070
addr_ws=${HOST_IP-127.0.0.1}:7000

sed -r -i "s/\/\/(\b[0-9]{1,3}\.){3}[0-9]{1,3}\b:[0-9]+/\/\/${addr_http}/g" site/static/js/main.06044d49.js
sed -r -i "s/\/\/(\b[0-9]{1,3}\.){3}[0-9]{1,3}\b:[0-9]+\/ws/\/\/${addr_ws}\/ws/g" site/static/js/main.06044d49.js

echo "build gochat.bin ..."
# CGO_CFLAGS="-g -O2 -Wno-return-local-addr" fix compile sqlite3 warning
CGO_CFLAGS="-g -O2 -Wno-return-local-addr" go build -o /go/src/gochat/bin/gochat.bin /go/src/gochat/main.go
echo "restart all ..."
supervisorctl restart all
echo "all Done."
echo "Beautiful ! Now, You can visit http://127.0.0.1:8080 , start the world."