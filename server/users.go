package server

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
)

type ServerUsers struct {
	socketHome  string
	socketName  string
	gid         int
	Names       map[string]*server
	waitStarted *sync.WaitGroup
}

func NewServerUsers(socketHome, socketName string) *ServerUsers {
	return &ServerUsers{
		socketHome:  socketHome,
		socketName:  socketName,
		gid:         -1,
		Names:       make(map[string]*server),
		waitStarted: &sync.WaitGroup{},
	}
}

func (s *ServerUsers) MakeFolder() error {
	_, err := os.Stat(s.socketHome)
	if err != nil && os.IsExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		err = os.Mkdir(s.socketHome, 0750)
		if err != nil {
			return err
		}
	}
	// FIXME chmod 750
	// FIXME set s.socketHome group to groupName
	return nil
}

func (s *ServerUsers) WithGroup(groupName string) (*ServerUsers, error) {
	g, err := user.LookupGroup(groupName)
	if err != nil {
		return nil, err
	}
	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return nil, err
	}
	s.gid = gid
	return s, nil
}

func (s *ServerUsers) AddUser(name string) (*server, error) {
	// verify the user exists on the system
	uzer, err := user.Lookup(name)
	if err != nil {
		return nil, err
	}

	listener, err := buildSocket(s.socketHome, s.socketName, uzer, s.gid)
	if err != nil {
		return nil, err
	}
	if listener == nil {
		panic("Socket can't be nil")
	}
	serv := NewServer(listener)
	s.Names[name] = serv
	return serv, nil
}

func (s *ServerUsers) Serve() {
	// FIXME use channels or Context to watch lifecycle of childrens
	wait := &sync.WaitGroup{}
	for _, server := range s.Names {
		wait.Add(1)
		go func() {
			server.Serve()
			wait.Done()
		}()
	}
	wait.Wait()
}

func (s *ServerUsers) Stop() {
	// FIXME stop
	w := &sync.WaitGroup{}
	fmt.Println("Names: ", s.Names)
	for name, server := range s.Names {
		w.Add(1)
		go func() {
			server.Stop()
			w.Done()
			fmt.Println("# server stopped: ", name)
		}()
	}
	w.Wait()
}

func uidgid(uzer *user.User) (uid int, guid int, err error) {
	// get uid user value as int
	uid, err = strconv.Atoi(uzer.Uid)
	if err != nil {
		return 0, 0, err
	}

	// get gid user value as int
	gid, err := strconv.Atoi(uzer.Gid)
	if err != nil {
		return 0, 0, err
	}

	return uid, gid, nil
}

func mkdirp(path string, perm os.FileMode) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, perm)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}
		err = os.Chmod(path, perm)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildSocket(home string, socketName string, uzer *user.User, withGid int) (*net.UnixListener, error) {
	uid, gid, err := uidgid(uzer)
	if err != nil {
		return nil, err
	}

	// socket dir
	sd := home + "/" + uzer.Username
	err = mkdirp(sd, 0700)
	if err != nil {
		return nil, err
	}

	sp := sd + "/" + socketName

	_, err = os.Stat(sp)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	err = os.Remove(sp)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") {
			// Who cares? it's just a cleanup
		} else {
			return nil, err
		}
	}

	listener, err := net.ListenUnix("unix", &net.UnixAddr{
		Name: sp,
		Net:  "unix",
	})
	if err != nil {
		return nil, err
	}

	err = os.Chmod(sp, 0600)
	if err != nil {
		return nil, err
	}

	// change socket ownsership to username
	if withGid != -1 {
		err = os.Chown(sd, uid, withGid)
		if err != nil {
			return nil, err
		}
	} else {
		err = os.Chown(sd, uid, gid)
		if err != nil {
			return nil, err
		}
	}

	if withGid != -1 {
		err = os.Chown(sp, uid, withGid)
		if err != nil {
			return nil, err
		}
	} else {
		err = os.Chown(sp, uid, gid)
		if err != nil {
			return nil, err
		}
	}

	if listener == nil {
		panic("conn is fucked up")
	}
	return listener, nil
}
