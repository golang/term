// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

func TestIsTerminalTerm(t *testing.T) {

	dir, err := ioutil.TempDir("", "TestIsTerminalTerm")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tty := filepath.Join(dir, "tty")
	err = unix.Mknod(tty, unix.S_IFCHR, int(unix.Mkdev(5, 0)))
	if err == unix.EPERM {
		t.Skip("no permission to create terminal file, skipping test")
	} else if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(tty)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if !term.IsTerminal(int(file.Fd())) {
		t.Fatalf("IsTerminal unexpectedly returned false for terminal file %s", file.Name())
	}
}
