## Key-Value Store
- Internally use `lmdb` for key-value storage.

#### Architecture
- An `HTTP Server` to send update request.
- An `TCP Server` to which clients can subscribe to keys.

`CGO_CFLAGS="-Wno-implicit-fallthrough" go run main.go -lmdbpath "/tmp/lmdbPath" -saddress "127.0.0.1:8200" -waddress "127.0.0.1:8300"`
