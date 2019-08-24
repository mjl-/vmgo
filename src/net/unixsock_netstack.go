// +build netstack

package net

import (
	"context"
	"os"
	"syscall"
)

func (c *UnixConn) readFrom(b []byte) (int, *UnixAddr, error) {
	return 0, nil, notyet("UnixConn.readFrom")
}

func (c *UnixConn) readMsg(b, oob []byte) (n, oobn, flags int, addr *UnixAddr, err error) {
	return 0, 0, 0, nil, notyet("UnixConn.readMsg")
}

func (c *UnixConn) writeTo(b []byte, addr *UnixAddr) (int, error) {
	return 0, notyet("UnixConn.writeTo")
}

func (c *UnixConn) writeMsg(b, oob []byte, addr *UnixAddr) (n, oobn int, err error) {
	return 0, 0, notyet("UnixConn.writeMsg")
}

func (ln *UnixListener) accept() (*UnixConn, error) {
	return nil, notyet("UnixListener.accept")
}

func (ln *UnixListener) close() error {
	return notyet("UnixListener.close")
}

func (ln *UnixListener) file() (*os.File, error) {
	return nil, syscall.ENOTSUP
}

func (sd *sysDialer) dialUnix(ctx context.Context, laddr, raddr *UnixAddr) (*UnixConn, error) {
	return nil, notyet("sysDialer.dialUnix")
}

func (sl *sysListener) listenUnix(ctx context.Context, laddr *UnixAddr) (*UnixListener, error) {
	return nil, notyet("sysListener.listenUnix")
}

func (sl *sysListener) listenUnixgram(ctx context.Context, laddr *UnixAddr) (*UnixConn, error) {
	return nil, notyet("sysListener.listenUnixgram")
}
