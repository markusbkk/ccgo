// Copyright 2021 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows
// +build !windows

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

// See the discussion at https://groups.google.com/g/golang-nuts/c/U1BLt5JM8F0/m/x9pQVq5RDwAJ
//
// Thanks to Brian Candler for fixing the code of this function.
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
	c := exec.Command(cmd, args...)
	c.Stdout = &b
	c.Stderr = &b
	c.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	// t0 := time.Now()
	err = c.Start()
	if err == nil {
		waitDone := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				unix.Kill(-c.Process.Pid, os.Kill.(syscall.Signal))
			case <-waitDone:
			}
		}()
		err = c.Wait()
		close(waitDone)
	}
	// if len(b.w.Bytes()) > 1e6 {
	// 	return nil, fmt.Errorf("output too big: %v", len(b.w.Bytes()))
	// }

	// if time.Since(t0) > time.Second {
	// 	return nil, fmt.Errorf("run too long: %v", time.Since(t0))
	// }

	return b.w.Bytes(), err
}
