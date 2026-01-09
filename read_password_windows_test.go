// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package term_test

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

func ExampleReadPassword_windows() {
	// If standard input is not bound to the terminal, a password can
	// still be read from it.

	path, err := windows.UTF16PtrFromString("CONIN$")
	if err != nil {
		panic(err)
	}
	var sa windows.SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	tty, err := windows.CreateFile(path, windows.GENERIC_READ|windows.GENERIC_WRITE, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, &sa, windows.OPEN_EXISTING, 0, 0)
	if err != nil {
		panic(err)
	}
	defer windows.CloseHandle(tty)

	fmt.Print("Enter your password: ")
	p, err := term.ReadPassword(int(tty))
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nYou entered '%s'\n", string(p))
}
