package server

import (
	"net"
	"os"
	"os/user"
	"strconv"
)

type ServerUsers struct {
	socketHome string
	gid        int
	Names      map[string]*server
}

func NewServerUsers(socketHome string) *ServerUsers {
	return &ServerUsers{
		socketHome: socketHome,
		gid:        -1,
		Names:      make(map[string]*server),
	}
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

func (s *ServerUsers) AddUser(name string) error {
	// verify the user exists on the system
	uzer, err := user.Lookup(name)
	if err != nil {
		return err
	}

	socket, err := buildSocket(s.socketHome, uzer)
	if err != nil {
		return err
	}

	s.Names[name] = NewServer(socket)
	return nil
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

func buildSocket(home string, uzer *user.User) (*net.UnixListener, error) {
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

	sp := sd + "/" + "restartctl.sock"

	_, err = os.Stat(sp)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	err = os.Remove(sp)
	if err == nil {
		return nil, err
	}

	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: sp, Net: "unix"})
	if err != nil {
		return nil, err
	}

	err = os.Chmod(sp, 0600)
	if err != nil {
		return nil, err
	}

	// change socket ownsership to username
	err = os.Chown(sd, uid, gid)
	if err != nil {
		return nil, err
	}

	err = os.Chown(sp, uid, gid)
	if err != nil {
		return nil, err
	}
	return l, nil

}
