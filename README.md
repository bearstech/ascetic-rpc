Ascetic RPC
===========

Ascetic protocol on a wire.

 - [x] Pascal string message
 - [x] Protobuf serialization for Request, Response, Body
 - [x] No auth, just UNIX socket with specific owner and group
 - [x] Server handles multiple socket, one per usage
 - [x] Graceful stop
 - [ ] Hot conf reload, with `kill -HUP`
 - [ ] Compression
 - [ ] Generating client proxy

Licence
-------
©2016 Mathieu Lecarme, 3 terms BSD Licence
