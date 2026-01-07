// +build linux
package posix

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	TCSANOW   = 0
	TCSADRAIN = 1
	TCSAFLUSH = 2
)

var (
	STDIN_FILENO  = unix.Stdin
	STDOUT_FILENO = unix.Stdout
	STDERR_FILENO = unix.Stderr
)

func Cfgetospeed(tio *unix.Termios) uint32 {
	return tio.Cflag & unix.CBAUD
}

func Cfgetispeed(tio *unix.Termios) uint32 {
	return Cfgetospeed(tio)
}

func Cfmakeraw(t *unix.Termios) {
	t.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	t.Oflag &^= unix.OPOST
	t.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	t.Cflag &^= unix.CSIZE | unix.PARENB
	t.Cflag |= unix.CS8
	t.Cc[unix.VMIN] = 1
	t.Cc[unix.VTIME] = 0
}

func Cfsetospeed(tio *unix.Termios, speed uint32) error {
	if speed&^unix.CBAUD != 0 {
		return unix.EINVAL
	}
	tio.Cflag &^= unix.CBAUD
	tio.Cflag |= speed
	return nil
}

func Cfsetispeed(tio *unix.Termios, speed uint32) error {
	if speed == 0 {
		return nil
	}
	return Cfsetospeed(tio, speed)
}

func Tcgetattr(fd int, value *unix.Termios) error {
	t, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if t != nil {
		*value = *t
	}
	return err
}

func Tcgetsid(fd int) (int, error) {
	return unix.IoctlGetInt(fd, unix.TIOCGSID)
}

func Tcsetattr(fd int, actions uint, value *unix.Termios) error {
	return unix.IoctlSetTermios(fd, unix.TCSETS+actions, value)
}

func Tcsendbreak(fd, dur int) error {
	return unix.IoctlSetInt(fd, unix.TCSBRK, 0)
}

func Tcflow(fd, action int) error {
	return unix.IoctlSetInt(fd, unix.TCXONC, action)
}

func Tcdrain(fd int) error {
	return unix.IoctlSetInt(fd, unix.TCSBRK, 1)
}

func Tcflush(fd, queue int) error {
	return unix.IoctlSetInt(fd, unix.TCFLSH, queue)
}

func Tcsetpgrp(fd, pgrp int) error {
	return unix.IoctlSetInt(fd, unix.TIOCSPGRP, pgrp)
}

func Openpty(tio *unix.Termios, ws *unix.Winsize) (m, s *os.File, err error) {
	defer func() {
		if err != nil {
			if m != nil {
				m.Close()
			}
			if s != nil {
				s.Close()
			}
		}
	}()

	m, err = os.OpenFile("/dev/ptmx", os.O_RDWR|unix.O_NOCTTY, 0644)
	if err != nil {
		return
	}

	err = Unlockpt(m)
	if err != nil {
		return
	}

	name, err := Ptsname(m)
	if err != nil {
		return
	}

	s, err = os.OpenFile(name, os.O_RDWR|unix.O_NOCTTY, 0644)
	if err != nil {
		return
	}

	sfd := int(s.Fd())
	if tio != nil {
		err = Tcsetattr(sfd, TCSANOW, tio)
		if err != nil {
			return
		}
	}

	if ws != nil {
		err = unix.IoctlSetWinsize(sfd, unix.TIOCSWINSZ, ws)
		if err != nil {
			return
		}
	}
	return
}

func Openpt(flags int) (*os.File, error) {
	return os.OpenFile("/dev/ptmx", flags, 0644)
}

func Unlockpt(f *os.File) error {
	unlock := 0
	fd := int(f.Fd())
	return unix.IoctlSetPointerInt(fd, unix.TIOCSPTLCK, unlock)
}

func Ptsname(f *os.File) (string, error) {
	fd := int(f.Fd())
	pty, err := unix.IoctlGetInt(fd, unix.TIOCGPTN)
	return fmt.Sprintf("/dev/pts/%d", pty), err
}

func Isatty(fd int) bool {
	var t unix.Termios
	return Tcgetattr(fd, &t) == nil
}

func Ttyname(fd int) (string, error) {
	if !Isatty(fd) {
		return "", syscall.ENOTTY
	}

	procname := fmt.Sprintf("/proc/self/fd/%d", fd)
	return os.Readlink(procname)
}
