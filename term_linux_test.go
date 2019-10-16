// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term_test

import (
	"os"
	"testing"

	"golang.org/x/term"
)

func TestIsTerminalTerm(t *testing.T) {
	file, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if !term.IsTerminal(int(file.Fd())) {
		t.Fatalf("IsTerminal unexpectedly returned false for terminal file %s", file.Name())
	}
}
