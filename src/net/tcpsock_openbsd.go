package net

import (
	"context"
	"io"
	"os"
	"syscall"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/network/ipv4"
)

func (ln *TCPListener) ok() bool { return ln != nil && ln.fd != nil }

func (ln *TCPListener) accept() (*TCPConn, error) {
	fd, err := ln.fd.accept()
	if err != nil {
		return nil, err
	}
	tc := newTCPConn(fd)
	if ln.lc.KeepAlive >= 0 {
		setKeepAlive(fd, true)
		ka := ln.lc.KeepAlive
		if ln.lc.KeepAlive == 0 {
			ka = defaultTCPKeepAlive
		}
		setKeepAlivePeriod(fd, ka)
	}
	return tc, nil
}

func (c *TCPConn) readFrom(r io.Reader) (int64, error) {
	return genericReadFrom(c, r)
}

func (ln *TCPListener) close() error {
	return ln.fd.Close()
}

func (ln *TCPListener) file() (*os.File, error) {
	return nil, syscall.ENOTSUP
}

func (sd *sysDialer) dialTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error) {
	if netstack == nil {
		return nil, errStack
	}

	lnsaddr := tcpip.FullAddress{}
	if laddr != nil && laddr.IP != nil {
		lnsaddr.Addr = tcpip.Address(laddr.IP.To4())
	}
	if laddr != nil {
		lnsaddr.Port = uint16(laddr.Port)
	}

	rnsaddr := tcpip.FullAddress{
		Addr: tcpip.Address(raddr.IP.To4()),
		Port: uint16(raddr.Port),
		// NIC
	}
	c, err := gonetDialContextTCP(ctx, netstack, rnsaddr, ipv4.ProtocolNumber)
	if err != nil {
		return nil, err
	}
	fd := &netFD{
		conn: c,
		net:  "tcp4", // xxx
		// xxx laddr, raddr
	}
	return newTCPConn(fd), nil
}

func (sl *sysListener) listenTCP(ctx context.Context, laddr *TCPAddr) (*TCPListener, error) {
	if netstack == nil {
		return nil, errStack
	}

	lnsaddr := tcpip.FullAddress{
		Port: uint16(laddr.Port),
	}
	if laddr.IP != nil {
		lnsaddr.Addr = tcpip.Address(laddr.IP.To4())
	}
	c, err := gonetNewListener(netstack, lnsaddr, ipv4.ProtocolNumber)
	if err != nil {
		return nil, err
	}
	nfd := &netFD{
		lconn: c,
		net:   "tcp4", // xxx
		laddr: laddr,
	}
	return &TCPListener{fd: nfd, lc: sl.ListenConfig}, nil
}
