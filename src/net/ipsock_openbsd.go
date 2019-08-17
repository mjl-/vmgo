// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"syscall"
)

func (p *ipStackCapabilities) probe() {
	p.ipv4Enabled = true
	p.ipv4MappedIPv6Enabled = false
	p.ipv6Enabled = false
}

func favoriteAddrFamily(network string, laddr, raddr sockaddr, mode string) (family int, ipv6only bool) {
	switch network[len(network)-1] {
	case '4':
		return syscall.AF_INET, false
	case '6':
		return syscall.AF_INET6, true
	}
	return syscall.AF_INET, false
}
