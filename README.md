## Key-Value Store
- Internally use `lmdb` for key-value storage.

#### Architecture
- An `HTTP Server` to send update request.
- An `TCP Server` to which clients can subscribe to keys.
