// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term_test

import (
	"io/ioutil"
	"os"
	"testing"

	"golang.org/x/term"
)

func TestIsTerminalTempFile(t *testing.T) {
	file, err := ioutil.TempFile("", "TestIsTerminalTempFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if term.IsTerminal(int(file.Fd())) {
		t.Fatalf("IsTerminal unexpectedly returned true for temporary file %s", file.Name())
	}
}
