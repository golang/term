// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term_test

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func ExampleReadPassword() {
	fmt.Print("Enter your password: ")
	p, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nYou entered '%s'\n", string(p))
}
