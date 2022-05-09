# a net.Conn compatible lib via QUIC

The initial code is cloned from [quic-conn](https://github.com/marten-seemann/quic-conn), but I refactored and updated codes to keepcompatible with latest https://github.com/lucas-clemente/quic-go.


## Run the example

Start listening for an incoming QUIC connection
```go
go run example/main.go -s
```
The server will echo every message received on the connection in uppercase.

Send a message on the QUIC connection:
```go
go run example/main.go -c
```
