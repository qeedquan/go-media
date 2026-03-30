// +build linux
package posix

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Passwd struct {
	Name   string
	Passwd string
	UID    int
	GID    int
	Gecos  string
	Dir    string
	Shell  string
}

func Getpwall() ([]Passwd, error) {
	var pw []Passwd
	var p Passwd

	f, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		t := strings.Split(s.Text(), ":")
		if len(t) >= 1 {
			p.Name = t[0]
		}
		if len(t) >= 2 {
			p.Passwd = t[1]
		}
		if len(t) >= 3 {
			p.UID, _ = strconv.Atoi(t[2])
		}
		if len(t) >= 4 {
			p.GID, _ = strconv.Atoi(t[3])
		}
		if len(t) >= 5 {
			p.Gecos = t[4]
		}
		if len(t) >= 6 {
			p.Dir = t[5]
		}
		if len(t) >= 7 {
			p.Shell = t[6]
		}
		pw = append(pw, p)
	}

	return pw, nil
}

func Getpwuid(uid int) (*Passwd, error) {
	pw, err := Getpwall()
	if err != nil {
		return nil, err
	}

	for i := range pw {
		if pw[i].UID == uid {
			return &pw[i], nil
		}
	}
	return nil, syscall.ENOENT
}
