# A QUIC Connection

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/marten-seemann/quic-conn)
[![Linux Build Status](https://img.shields.io/travis/marten-seemann/quic-conn/master.svg?style=flat-square&label=linux+build)](https://travis-ci.org/marten-seemann/quic-conn)
[![Code Coverage](https://img.shields.io/codecov/c/github/marten-seemann/quic-conn/master.svg?style=flat-square)](https://codecov.io/gh/marten-seemann/quic-conn/)

At the moment, this project is intended to figure out the right API exposed by the [quic package in quic-go](https://github.com/lucas-clemente/quic-go).

When fully implemented, a QUIC connection can be used as a replacement for an encrypted TCP connection. It provides a single ordered byte-stream abstraction, with the main benefit of being able to perform connection migration.

## Usage of the example

Start listening for an incoming QUIC connection
```go
go run example/main.go -s
```
The server will echo every message received on the connection in uppercase.

Send a message on the QUIC connection:
```go
go run example/main.go -c
```
