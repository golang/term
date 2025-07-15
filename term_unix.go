// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package term

import (
	"context"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

type state struct {
	termios unix.Termios
}

func isTerminal(fd int) bool {
	_, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	return err == nil
}

func makeRaw(fd int) (*State, error) {
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, err
	}

	oldState := State{state{termios: *termios}}

	// This attempts to replicate the behaviour documented for cfmakeraw in
	// the termios(3) manpage.
	termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, termios); err != nil {
		return nil, err
	}

	return &oldState, nil
}

func getState(fd int) (*State, error) {
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, err
	}

	return &State{state{termios: *termios}}, nil
}

func restore(fd int, state *State) error {
	return unix.IoctlSetTermios(fd, ioctlWriteTermios, &state.termios)
}

func getSize(fd int) (width, height int, err error) {
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Col), int(ws.Row), nil
}

// passwordReader is an io.Reader that reads from a specific file descriptor.
type passwordReader int

func (r passwordReader) Read(buf []byte) (int, error) {
	return unix.Read(int(r), buf)
}

func readPassword(fd int) ([]byte, error) {
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, err
	}

	newState := *termios
	newState.Lflag &^= unix.ECHO
	newState.Lflag |= unix.ICANON | unix.ISIG
	newState.Iflag |= unix.ICRNL
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		return nil, err
	}

	defer unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)

	return readPasswordLine(passwordReader(fd))
}

func readPasswordWithContext(fd int, ctx context.Context) ([]byte, error) {
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	nonblocking := false
	if err != nil {
		return nil, err
	}
	newState := *termios
	newState.Lflag &^= unix.ECHO
	newState.Lflag |= unix.ICANON | unix.ISIG
	newState.Iflag |= unix.ICRNL

	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		return nil, err
	}
	defer func() {
		if nonblocking {
			unix.SetNonblock(fd, false)
		}
		unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)
	}()

	// Set nonblocking IO
	if err := unix.SetNonblock(fd, true); err != nil {
		return nil, err
	}
	nonblocking = true

	var ret []byte
	var buf [1]byte
	for {
		if ctx.Err() != nil {
			return ret, ctx.Err()
		}
		n, err := unix.Read(fd, buf[:])
		if err != nil {
			// Check for nonblocking error
			if serr, ok := err.(syscall.Errno); ok {
				if serr == 11 {
					// Add (hopefully not noticable) latency to prevent CPU hogging
					time.Sleep(50 * time.Millisecond)
					continue
				}
			}
			return ret, err
		}
		if n > 0 {
			switch buf[0] {
			case '\b':
				if len(ret) > 0 {
					ret = ret[:len(ret)-1]
				}
			case '\n':
				return ret, nil
			default:
				ret = append(ret, buf[0])
			}
			continue
		}
	}
}
