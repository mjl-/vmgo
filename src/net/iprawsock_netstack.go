// +build netstack

package net

import (
	"context"
)

func (c *IPConn) readFrom(b []byte) (int, *IPAddr, error) {
	return 0, nil, notyet("IPConn.readFrom")
}

func (c *IPConn) readMsg(b, oob []byte) (n, oobn, flags int, addr *IPAddr, err error) {
	return 0, 0, 0, nil, notyet("IPConn.readMsg")
}

func (c *IPConn) writeTo(b []byte, addr *IPAddr) (int, error) {
	return 0, notyet("IPConn.writeTo")
}

func (c *IPConn) writeMsg(b, oob []byte, addr *IPAddr) (n, oobn int, err error) {
	return 0, 0, notyet("IPConn.writeMsg")
}

func (sd *sysDialer) dialIP(ctx context.Context, laddr, raddr *IPAddr) (*IPConn, error) {
	return nil, notyet("sysDialer.dialIP")
}

func (sl *sysListener) listenIP(ctx context.Context, laddr *IPAddr) (*IPConn, error) {
	return nil, notyet("sysListener.listenIP")
}
