// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"time"

	"github.com/google/netstack/tcpip"
)

func setKeepAlivePeriod(fd *netFD, d time.Duration) error {
	err := fd.conn.ep.SetSockOpt(tcpip.KeepaliveIdleOption(d))
	if err != nil {
		return wrapNetstackError(err)
	}
	return wrapNetstackError(fd.conn.ep.SetSockOpt(tcpip.KeepaliveIntervalOption(d)))
}

func setNoDelay(fd *netFD, noDelay bool) error {
	v := 1
	if noDelay {
		v = 0
	}
	return wrapNetstackError(fd.conn.ep.SetSockOpt(tcpip.DelayOption(v)))
}
