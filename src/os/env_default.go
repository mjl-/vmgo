// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows
// +build !solo5hvt

package os

import "syscall"

func environForSysProcAttr(sys *syscall.SysProcAttr) ([]string, error) {
	return Environ(), nil
}
