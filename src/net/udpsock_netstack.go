// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build netstack

package net

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/network/ipv4"
	"github.com/google/netstack/tcpip/network/ipv6"
)

// stub
func sockaddrToUDP(sa syscall.Sockaddr) Addr {
	return nil
}

// UDPAddr represents the address of a UDP end point.
type UDPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

// Network returns the address's network name, "udp".
func (a *UDPAddr) Network() string { return "udp" }

func (a *UDPAddr) String() string {
	if a == nil {
		return "<nil>"
	}
	ip := ipEmptyString(a.IP)
	if a.Zone != "" {
		return JoinHostPort(ip+"%"+a.Zone, itoa(a.Port))
	}
	return JoinHostPort(ip, itoa(a.Port))
}

func (a *UDPAddr) isWildcard() bool {
	if a == nil || a.IP == nil {
		return true
	}
	return a.IP.IsUnspecified()
}

func (a *UDPAddr) opAddr() Addr {
	if a == nil {
		return nil
	}
	return a
}

// ResolveUDPAddr returns an address of UDP end point.
//
// The network must be a UDP network name.
//
// If the host in the address parameter is not a literal IP address or
// the port is not a literal port number, ResolveUDPAddr resolves the
// address to an address of UDP end point.
// Otherwise, it parses the address as a pair of literal IP address
// and port number.
// The address parameter can use a host name, but this is not
// recommended, because it will return at most one of the host name's
// IP addresses.
//
// See func Dial for a description of the network and address
// parameters.
func ResolveUDPAddr(network, address string) (*UDPAddr, error) {
	switch network {
	case "udp", "udp4", "udp6":
	case "": // a hint wildcard for Go 1.0 undocumented behavior
		network = "udp"
	default:
		return nil, UnknownNetworkError(network)
	}
	addrs, err := DefaultResolver.internetAddrList(context.Background(), network, address)
	if err != nil {
		return nil, err
	}
	return addrs.forResolve(network, address).(*UDPAddr), nil
}

// UDPConn is the implementation of the Conn and PacketConn interfaces
// for UDP network connections.
type UDPConn struct {
	nsconn *gonetPacketConn
}

func (c *UDPConn) Close() error {
	return c.nsconn.Close()
}

func (c *UDPConn) File() (*os.File, error) {
	return nil, syscall.ENOTSUP
}

func (c *UDPConn) LocalAddr() Addr {
	return c.nsconn.LocalAddr()
}

func (c *UDPConn) Read(b []byte) (int, error) {
	return c.nsconn.Read(b)
}

func (c *UDPConn) ReadFrom(b []byte) (int, Addr, error) {
	return c.nsconn.ReadFrom(b)
}

func (c *UDPConn) ReadFromUDP(b []byte) (int, *UDPAddr, error) {
	n, addr, err := c.nsconn.ReadFrom(b)
	var udpAddr *UDPAddr
	if addr != nil {
		ipStr, portStr, err := SplitHostPort(addr.String())
		if err != nil {
			return n, nil, err // todo: return an OpError
		}
		port, _ := parsePort(portStr)
		udpAddr = &UDPAddr{
			IP:   ParseIP(ipStr),
			Port: port,
			// todo for ipv6: Zone
		}
	}
	return n, udpAddr, err
}

// ReadMsgUDP reads a message from c, copying the payload into b and
// the associated out-of-band data into oob. It returns the number of
// bytes copied into b, the number of bytes copied into oob, the flags
// that were set on the message and the source address of the message.
//
// The packages golang.org/x/net/ipv4 and golang.org/x/net/ipv6 can be
// used to manipulate IP-level socket options in oob.
func (c *UDPConn) ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *UDPAddr, err error) {
	if len(oob) != 0 {
		return 0, 0, 0, nil, syscall.ENOTSUP
	}
	n, addr, err = c.ReadFromUDP(b)
	return
}

func (c *UDPConn) RemoteAddr() Addr {
	return c.nsconn.RemoteAddr()
}

// SyscallConn returns a raw network connection.
// Not implemented.
func (c *UDPConn) SyscallConn() (syscall.RawConn, error) {
	return nil, syscall.ENOTSUP
}

func (c *UDPConn) Write(b []byte) (int, error) {
	return c.nsconn.Write(b)
}

func (c *UDPConn) WriteMsgUDP(b, oob []byte, addr *UDPAddr) (n, oobn int, err error) {
	if len(oob) != 0 {
		return 0, 0, syscall.ENOTSUP
	}
	n, err = c.nsconn.WriteTo(b, addr)
	return
}

func (c *UDPConn) WriteTo(b []byte, addr Addr) (int, error) {
	return c.nsconn.WriteTo(b, addr)
}

func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error) {
	return c.nsconn.WriteTo(b, addr)
}

func (c *UDPConn) SetReadBuffer(bytes int) error {
	return syscall.ENOTSUP
}

func (c *UDPConn) SetWriteBuffer(bytes int) error {
	return syscall.ENOTSUP
}

func (c *UDPConn) SetDeadline(t time.Time) error {
	return c.nsconn.SetDeadline(t)
}

func (c *UDPConn) SetReadDeadline(t time.Time) error {
	return c.nsconn.SetReadDeadline(t)
}

func (c *UDPConn) SetWriteDeadline(t time.Time) error {
	return c.nsconn.SetWriteDeadline(t)
}

func udpAddrToFull(a *UDPAddr) (tcpip.NetworkProtocolNumber, *tcpip.FullAddress) {
	proto := ipv6.ProtocolNumber
	if a == nil {
		return proto, nil
	}
	fa := &tcpip.FullAddress{
		Port: uint16(a.Port),
	}
	if a.IP != nil {
		if a.IP.To4() != nil {
			fa.Addr = tcpip.Address(a.IP.To4())
			proto = ipv4.ProtocolNumber
		} else {
			fa.Addr = tcpip.Address(a.IP)
		}
	}
	return proto, fa
}

// DialUDP acts like Dial for UDP networks.
//
// The network must be a UDP network name; see func Dial for details.
//
// If laddr is nil, a local address is automatically chosen.
// If the IP field of raddr is nil or an unspecified IP address, the
// local system is assumed.
func DialUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error) {
	if netstack == nil {
		return nil, errStack
	}

	switch network {
	case "udp", "udp4", "udp6":
	default:
		return nil, &OpError{Op: "dial", Net: network, Source: laddr.opAddr(), Addr: raddr.opAddr(), Err: UnknownNetworkError(network)}
	}
	if raddr == nil {
		return nil, &OpError{Op: "dial", Net: network, Source: laddr.opAddr(), Addr: nil, Err: errMissingAddress}
	}

	_, lnsaddr := udpAddrToFull(laddr)
	proto, rnsaddr := udpAddrToFull(raddr)
	c, err := gonetDialUDP(netstack, lnsaddr, rnsaddr, proto)
	if err != nil {
		return nil, err
	}
	uc := &UDPConn{c}
	return uc, nil
}

// ListenUDP acts like ListenPacket for UDP networks.
//
// The network must be a UDP network name; see func Dial for details.
//
// If the IP field of laddr is nil or an unspecified IP address,
// ListenUDP listens on all available IP addresses of the local system
// except multicast IP addresses.
// If the Port field of laddr is 0, a port number is automatically
// chosen.
func ListenUDP(network string, laddr *UDPAddr) (*UDPConn, error) {
	if netstack == nil {
		return nil, errStack
	}

	switch network {
	case "udp", "udp4", "udp6":
	default:
		return nil, &OpError{Op: "listen", Net: network, Source: nil, Addr: laddr.opAddr(), Err: UnknownNetworkError(network)}
	}
	if laddr == nil {
		laddr = &UDPAddr{}
	}
	proto, lnsaddr := udpAddrToFull(laddr)
	c, err := gonetDialUDP(netstack, lnsaddr, nil, proto)
	if err != nil {
		return nil, err
	}
	uc := &UDPConn{c}
	return uc, nil
}

// ListenMulticastUDP acts like ListenPacket for UDP networks but
// takes a group address on a specific network interface.
//
// The network must be a UDP network name; see func Dial for details.
//
// ListenMulticastUDP listens on all available IP addresses of the
// local system including the group, multicast IP address.
// If ifi is nil, ListenMulticastUDP uses the system-assigned
// multicast interface, although this is not recommended because the
// assignment depends on platforms and sometimes it might require
// routing configuration.
// If the Port field of gaddr is 0, a port number is automatically
// chosen.
//
// ListenMulticastUDP is just for convenience of simple, small
// applications. There are golang.org/x/net/ipv4 and
// golang.org/x/net/ipv6 packages for general purpose uses.
func ListenMulticastUDP(network string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error) {
	return nil, notyet("ListenMulticastUDP")
}

func (sd *sysDialer) dialUDP(ctx context.Context, laddr, raddr *UDPAddr) (*UDPConn, error) {
	if netstack == nil {
		return nil, errStack
	}

	_, lnsaddr := udpAddrToFull(laddr)
	proto, rnsaddr := udpAddrToFull(raddr)
	c, err := gonetDialUDP(netstack, lnsaddr, rnsaddr, proto)
	if err != nil {
		return nil, err
	}
	uc := &UDPConn{c}
	return uc, nil
}

func (sl *sysListener) listenUDP(ctx context.Context, laddr *UDPAddr) (*UDPConn, error) {
	if netstack == nil {
		return nil, errStack
	}

	proto, lnsaddr := udpAddrToFull(laddr)
	c, err := gonetDialUDP(netstack, lnsaddr, nil, proto)
	if err != nil {
		return nil, err
	}
	uc := &UDPConn{c}
	return uc, nil
}

func (sl *sysListener) listenMulticastUDP(ctx context.Context, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error) {
	return nil, notyet("sysListener.listenMulticastUDP")
}
