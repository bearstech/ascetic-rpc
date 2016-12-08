package server

import (
	"errors"
	"os/user"
	"testing"

	"github.com/bearstech/ascetic-rpc/client"
	"github.com/bearstech/ascetic-rpc/model"
)

// hello is declared in server_test.go

func TestUsersHello(t *testing.T) {
	me, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	servers := NewServerUsers("/tmp/test", "ascetic.sock")

	err = servers.MakeFolder()
	if err != nil {
		t.Fatal(err)
	}
	myserver, err := servers.AddUser(me.Username)
	if err != nil {
		t.Fatal(err)
	}
	myserver.Register("hello", hello)
	servers.Serve()

	c, err := client.NewClientUnix("/tmp/test/" + me.Username + "/ascetic.sock")
	if err != nil {
		t.Fatal(err)
	}

	hello := model.Hello{Name: "Charlie"}
	var world model.World
	err = c.Do("hello", &hello, &world)
	if err != nil {
		t.Fatal(err)
	}
	if world.Message != "Hello Charlie♥️" {
		t.Error(errors.New("Bad message: " + world.Message))
	}

	servers.Stop()
}
