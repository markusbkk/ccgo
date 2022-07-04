// Copyright 2021 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"context"
	"fmt"
	"os/exec"
)

func shell(echo bool, cmd string, args ...string) ([]byte, error) {
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		return nil, err
	}

	wd, err := absCwd()
	if err != nil {
		return nil, err
	}

	if echo {
		fmt.Printf("execute %s %q in %s\n", cmd, args, wd)
	}
	var b echoWriter
	b.silent = !echo
	ctx, cancel := context.WithTimeout(context.Background(), *oShellTime)

	defer cancel()

	c := exec.CommandContext(ctx, cmd, args...)
	c.Stdout = &b
	c.Stderr = &b
	err = c.Run()
	return b.w.Bytes(), err
}
