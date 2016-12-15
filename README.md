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

Documentation
-------------

Here is the "hello world" example.

Describe your input/output, or request/response, with protobuf:

```protbuf
syntax = "proto3";
package model;

message Hello {
    string Name = 1;
}

message World {
    string Message = 1;
}
```

The function:

```golang
import (
    "fmt"
    "github.com/bearstech/ascetic-rpc/model"
)

func hello(req *model.Request) (*model.Response, error) {
    var hello model.Hello
    err := req.GetBody(&hello)
    if err != nil {
        panic(err)
    }
    world := model.World{
        Message: fmt.Sprintf("Hello %s♥️", hello.Name),
    }
    res, err := model.NewOKResponse(&world)
    if err != nil {
        return nil, err
    }
    return res, nil
}
```

The server part:

```golang
import (
    "github.com/bearstech/ascetic-rpc/server"
)

func main() {
    s, err := server.NewServerUnix("/tmp/example.sock")
    if err != nil {
        panic(err)
    }
    defer s.Stop()
    s.register("hello", hello)
    s.Serve()
}
```

The Client part:

```golang
import (
    "fmt"
    "github.com/bearstech/ascetic-rpc/client"
    "github.com/bearstech/ascetic-rpc/model"
)

func main() {
    c, err := client.NewClientUnix("/tmp/example.sock")
    if err != nil {
        panic(err)
    }
    var world model.World
    hello := mode.Hello{Name: "Bob"}
    err := c.Do("hello", &hello, &world)
    if err != nil {
        panic(err)
    }
    fmt.Println(world.Message)
}
```


Licence
-------
©2016 Mathieu Lecarme, 3 terms BSD Licence
