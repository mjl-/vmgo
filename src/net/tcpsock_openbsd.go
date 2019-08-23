package net

import (
	"context"
	"io"
	"os"
	"syscall"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/network/ipv4"
	"github.com/google/netstack/tcpip/network/ipv6"
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

// ipToAddress returns an address and whether it is ipv4-only.
func ipToAddress(ip IP) (tcpip.NetworkProtocolNumber, tcpip.Address) {
	v4 := ip.To4()
	if v4 != nil {
		return ipv4.ProtocolNumber, tcpip.Address(v4)
	}
	return ipv6.ProtocolNumber, tcpip.Address(ip)
}

func (sd *sysDialer) dialTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error) {
	if netstack == nil {
		return nil, errStack
	}

	lnsaddr := tcpip.FullAddress{}
	if laddr != nil && laddr.IP != nil {
		_, addr := ipToAddress(laddr.IP)
		lnsaddr.Addr = addr
	}
	if laddr != nil {
		lnsaddr.Port = uint16(laddr.Port)
	}

	protoNum, raddress := ipToAddress(raddr.IP)
	rnsaddr := tcpip.FullAddress{
		Addr: raddress,
		Port: uint16(raddr.Port),
		// NIC
	}
	c, err := gonetDialContextTCP(ctx, netstack, rnsaddr, protoNum)
	if err != nil {
		return nil, err
	}
	net := "tcp6"
	if protoNum == ipv4.ProtocolNumber {
		net = "tcp4"
	}
	fd := &netFD{
		conn:  c,
		net:   net,
		laddr: laddr,
		raddr: raddr,
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
	protoNum := ipv6.ProtocolNumber
	if laddr.IP != nil {
		protoNum, lnsaddr.Addr = ipToAddress(laddr.IP)
	}
	c, err := gonetNewListener(netstack, lnsaddr, protoNum)
	if err != nil {
		return nil, err
	}
	net := "tcp6"
	if protoNum == ipv4.ProtocolNumber {
		net = "tcp4"
	}
	nfd := &netFD{
		lconn: c,
		net:   net,
		laddr: laddr,
	}
	return &TCPListener{fd: nfd, lc: sl.ListenConfig}, nil
}
