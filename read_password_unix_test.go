// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package term_test

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func ExampleReadPassword_unix() {
	// If standard input is not bound to the terminal, a password can
	// still be read from it.
	tty, err := os.Open("/dev/tty")
	if err != nil {
		panic(err)
	}
	defer tty.Close()
	fmt.Print("Enter your password: ")
	p, err := term.ReadPassword(int(tty.Fd()))
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nYou entered '%s'\n", string(p))
}
