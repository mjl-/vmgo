package net

import (
	"os"
	"syscall"
	"time"
)

type netFD struct {
	// one of these should be set
	conn  *gonetConn
	lconn *gonetListener

	net   string
	laddr Addr
	raddr Addr
}

func (fd *netFD) ok() bool {
	return fd != nil && (fd.conn != nil || fd.lconn != nil)
}

func (fd *netFD) Read(p []byte) (n int, err error) {
	return fd.conn.Read(p)
}

func (fd *netFD) Write(p []byte) (nn int, err error) {
	return fd.conn.Write(p)
}

func (fd *netFD) Close() error {
	return fd.conn.Close()
}

func (fd *netFD) closeRead() error {
	return fd.conn.CloseRead()
}

func (fd *netFD) closeWrite() error {
	return fd.conn.CloseWrite()
}

func (fd *netFD) SetDeadline(t time.Time) error {
	return fd.conn.SetDeadline(t)
}

func (fd *netFD) SetReadDeadline(t time.Time) error {
	return fd.conn.SetReadDeadline(t)
}

func (fd *netFD) SetWriteDeadline(t time.Time) error {
	return fd.conn.SetWriteDeadline(t)
}

func (fd *netFD) accept() (netfd *netFD, err error) {
	c, err := fd.lconn.Accept()
	if err != nil {
		return nil, err
	}

	nfd := &netFD{
		conn:  c,
		net:   fd.net,
		laddr: c.LocalAddr(),
		raddr: c.RemoteAddr(),
	}
	return nfd, nil
}

func sysSocket(family, sotype, proto int) (int, error) {
	return 0, syscall.ENOSYS
}

func (fd *netFD) readFrom(p []byte) (n int, sa syscall.Sockaddr, err error) {
	return 0, nil, syscall.ENOSYS
}

func (fd *netFD) readMsg(p []byte, oob []byte) (n, oobn, flags int, sa syscall.Sockaddr, err error) {
	return 0, 0, 0, nil, syscall.ENOSYS
}

func (fd *netFD) writeTo(p []byte, sa syscall.Sockaddr) (n int, err error) {
	return 0, syscall.ENOSYS
}

func (fd *netFD) writeMsg(p []byte, oob []byte, sa syscall.Sockaddr) (n int, oobn int, err error) {
	return 0, 0, syscall.ENOSYS
}

func (fd *netFD) dup() (f *os.File, err error) {
	return nil, syscall.ENOSYS
}
