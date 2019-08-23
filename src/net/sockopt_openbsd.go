// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"syscall"

	"github.com/google/netstack/tcpip"
)

/*
import (
	"internal/bytealg"
	"runtime"
	"syscall"
)

// Boolean to int.
func boolint(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ipv4AddrToInterface(ip IP) (*Interface, error) {
	ift, err := Interfaces()
	if err != nil {
		return nil, err
	}
	for _, ifi := range ift {
		ifat, err := ifi.Addrs()
		if err != nil {
			return nil, err
		}
		for _, ifa := range ifat {
			switch v := ifa.(type) {
			case *IPAddr:
				if ip.Equal(v.IP) {
					return &ifi, nil
				}
			case *IPNet:
				if ip.Equal(v.IP) {
					return &ifi, nil
				}
			}
		}
	}
	if ip.Equal(IPv4zero) {
		return nil, nil
	}
	return nil, errNoSuchInterface
}

func interfaceToIPv4Addr(ifi *Interface) (IP, error) {
	if ifi == nil {
		return IPv4zero, nil
	}
	ifat, err := ifi.Addrs()
	if err != nil {
		return nil, err
	}
	for _, ifa := range ifat {
		switch v := ifa.(type) {
		case *IPAddr:
			if v.IP.To4() != nil {
				return v.IP, nil
			}
		case *IPNet:
			if v.IP.To4() != nil {
				return v.IP, nil
			}
		}
	}
	return nil, errNoSuchInterface
}

func setIPv4MreqToInterface(mreq *syscall.IPMreq, ifi *Interface) error {
	if ifi == nil {
		return nil
	}
	ifat, err := ifi.Addrs()
	if err != nil {
		return err
	}
	for _, ifa := range ifat {
		switch v := ifa.(type) {
		case *IPAddr:
			if a := v.IP.To4(); a != nil {
				copy(mreq.Interface[:], a)
				goto done
			}
		case *IPNet:
			if a := v.IP.To4(); a != nil {
				copy(mreq.Interface[:], a)
				goto done
			}
		}
	}
done:
	if bytealg.Equal(mreq.Multiaddr[:], IPv4zero.To4()) {
		return errNoSuchMulticastInterface
	}
	return nil
}
*/

func setReadBuffer(fd *netFD, bytes int) error {
	return wrapNetstackError(fd.conn.ep.SetSockOpt(tcpip.ReceiveBufferSizeOption(bytes)))
}

func setWriteBuffer(fd *netFD, bytes int) error {
	return wrapNetstackError(fd.conn.ep.SetSockOpt(tcpip.SendBufferSizeOption(bytes)))
}

func setKeepAlive(fd *netFD, keepalive bool) error {
	v := 0
	if keepalive {
		v = 1
	}
	return wrapNetstackError(fd.conn.ep.SetSockOpt(tcpip.KeepaliveEnabledOption(v)))
}

func setLinger(fd *netFD, sec int) error {
	// does not seem supported by netstack. if close would abort (not flushing pending
	// writes), we could do a flush in a goroutine during close with so_linger set. but
	// gonet doesn't have a flush operation.
	return syscall.ENOTSUP
}
