package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
)

var uid_arg string
var gid_arg string

func main() {
	err := realMain()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func get_uid(user_name string) error {
	u, er := user.Lookup(user_name)
	if er != nil {
		return er
	}
	uid_arg = u.Uid
	return er
}

func get_gid(group_name string) error {
	g, er := user.LookupGroup(group_name)
	if er != nil {
		return er
	}
	gid_arg = g.Gid
	return er
}

func realMain() error {
	var err error
	if len(os.Args) != 6 {
		return nil
	}

	// savetofile <mode> <filepath> <username> <groupname> <data>
	mode := os.Args[1]
	path := os.Args[2]
	user_name := os.Args[3]
	group_name := os.Args[4]
	data := os.Args[5]

	err = get_uid(user_name)
	if err != nil {
		return err
	}

	err = get_gid(group_name)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(uid_arg)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(gid_arg)
	if err != nil {
		return err
	}

	switch mode {
	case "append":
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0700)
		if err != nil {
			return err
		}

		defer f.Close()

		if _, err = f.WriteString(data); err != nil {
			return err
		}
	case "append-nl":
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0700)
		if err != nil {
			return err
		}

		defer f.Close()

		if _, err = f.WriteString(data); err != nil {
			return err
		}

		if _, err = f.WriteString("\n"); err != nil {
			return err
		}
	case "create-nl":
		err := ioutil.WriteFile(path, append([]byte(data), []byte("\n")...), 0700)
		if err != nil {
			return err
		}
	default: // "create"
		err := ioutil.WriteFile(path, []byte(data), 0700)
		if err != nil {
			return err
		}
	}

	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}

	err = os.Chmod(path, 0640)
	if err != nil {
		return err
	}

	return nil
}

