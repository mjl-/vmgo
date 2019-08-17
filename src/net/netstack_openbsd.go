package net

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/link/fdbased"
	"github.com/google/netstack/tcpip/link/sniffer"
	"github.com/google/netstack/tcpip/network/arp"
	"github.com/google/netstack/tcpip/network/ipv4"
	"github.com/google/netstack/tcpip/network/ipv6"
	"github.com/google/netstack/tcpip/stack"
	"github.com/google/netstack/tcpip/transport/tcp"
	"github.com/google/netstack/tcpip/transport/udp"
)

var (
	netstack            *stack.Stack
	netstackNameservers []string
	errStack            = errors.New("no network stack")
)

func init() {
	gonet := os.Getenv("GONET")
	if gonet == "" {
		// no stack configured
		return
	}

	rand.Seed(time.Now().UnixNano())

	fail := func(format string, args ...interface{}) {
		panic(fmt.Sprintf("netstack init: "+format, args...))
	}

	// Example of what we are trying to parse:
	//
	//	verbose
	// 	nic id=1 ether=12:12:12:23:23:23 mtu=1500 fd=3 dev=/dev/tap0 sniff=true
	// 	ip nic=1 addr=1.2.3.4/24 addr=2.3.4.5/32 addr=ff:aa.../64
	// 	route nic=1 ipnet=ip/mask gw=ip...
	// 	dns ip=1.2.3.4 ip=8.8.8.8

	parseArgs := func(pairs []string) [][2]string {
		args := make([][2]string, len(pairs))
		for i, p := range pairs {
			t := strings.SplitN(p, "=", 2)
			if len(t) != 2 {
				fail("bad pair %q", p)
			}
			args[i] = [...]string{t[0], t[1]}
		}
		return args
	}

	type nicip struct {
		nic tcpip.NICID
		ip  IP
	}
	var (
		verbose                     bool
		nics                        = map[tcpip.NICID]tcpip.LinkEndpointID{}
		ethers                      = map[tcpip.NICID]struct{}{}
		ipaddrs                     []nicip
		routes                      []tcpip.Route
		haveIP4, haveIP6, haveEther bool
	)

	xparseNIC := func(s string) tcpip.NICID {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			fail("invalid nic id %q: %v", s, err)
		}
		if v == 0 {
			fail("invalid nic id 0")
		}
		return tcpip.NICID(v)
	}

	parseNIC := func(args [][2]string) {
		var (
			nic   tcpip.NICID
			ether HardwareAddr
			mtu   int
			dev   string
			fd    int
			sniff bool
		)
		for _, p := range args {
			k, v := p[0], p[1]
			switch k {
			case "id":
				nic = xparseNIC(v)
			case "ether":
				if ether != nil {
					fail("duplicate key ether")
				}
				mac, err := ParseMAC(v)
				if err != nil {
					fail("parsing ether %q: %v", v, err)
				}
				if len(mac) != 6 {
					fail("invalid ether length, must be six bytes, saw %q of %d", v, len(mac))
				}
				ether = mac
				haveEther = true
			case "mtu":
				var err error
				mtu, err = strconv.Atoi(v)
				if err != nil {
					fail("parsing mtu %q: %v", v, err)
				}
			case "fd":
				if v != "3" {
					fail("only fd=3 supported")
				}
				fd = 3
			case "dev":
				dev = v
				fail("dev not yet")
			case "sniff":
				switch v {
				case "true":
					sniff = true
				case "false":
					sniff = false
				default:
					fail("bad sniff %q", v)
				}
			default:
				fail("bad route keyword %q", k)
			}
		}
		if nic == 0 {
			fail("nic: missing id")
		}
		if _, ok := nics[nic]; ok {
			fail("nic: duplicate nic %d", nic)
		}
		if mtu == 0 {
			fail("nic: missing mtu")
		}
		if fd == 0 && dev == "" {
			fail("nic: missing fd or dev")
		}
		if fd != 0 && dev != "" {
			fail("nic: cannot have both fd and dev")
		}

		fdOpts := &fdbased.Options{
			FDs:                []int{int(os.Stdtuntap.Fd())},
			MTU:                uint32(mtu),
			PacketDispatchMode: fdbased.Readv,
			EthernetHeader:     ether != nil,
			Address:            tcpip.LinkAddress(ether),
		}
		linkID, err := fdbased.New(fdOpts)
		if err != nil {
			fail("making new fd-based link: %v", err)
		}
		if sniff {
			linkID = sniffer.New(linkID)
		}
		nics[nic] = linkID
		if ether != nil {
			ethers[nic] = struct{}{}
		}
		if verbose {
			log.Printf("netstack: adding nic id=%d ether=%s mtu=%d fd=%d dev=%s sniff=%v", nic, ether, mtu, fd, dev, sniff)
		}
	}

	parseIP := func(args [][2]string) {
		var ips []IP
		var nic tcpip.NICID
		for _, p := range args {
			k, v := p[0], p[1]
			switch k {
			case "nic":
				nic = xparseNIC(v)
			case "addr":
				ip := ParseIP(v)
				if ip == nil {
					fail("bad addr value %q", v)
				}
				ips = append(ips, ip)
				if ip.To4() != nil {
					haveIP4 = true
				} else {
					haveIP6 = true
				}
			default:
				fail("bad ip keyword %q", k)
			}
		}
		if nic == 0 {
			fail("missing nic in ip statement")
		}
		if _, ok := nics[nic]; !ok {
			fail("unknown nic %d in ip statement, define it before use", nic)
		}
		if len(ips) == 0 {
			fail("no addresses in ip statement")
		}
		for _, ip := range ips {
			if ip.To4() != nil {
				ip = ip.To4()
			}
			a := nicip{nic, ip}
			if verbose {
				log.Printf("netstack: adding ip nic=%d addr=%s", a.nic, a.ip)
			}
			ipaddrs = append(ipaddrs, a)
		}
	}

	parseRoute := func(args [][2]string) {
		var (
			nic    tcpip.NICID
			ipnet  *IPNet
			gw     IP
			haveGW bool
		)
		for _, p := range args {
			k, v := p[0], p[1]
			switch k {
			case "ipnet":
				var err error
				_, ipnet, err = ParseCIDR(v)
				if err != nil {
					fail("bad ipnet %q in route statement: %v", v, err)
				}
			case "gw":
				if v == "" {
					haveGW = true
					continue
				}
				gw = ParseIP(v)
				if gw == nil {
					fail("bad gw %q in route statement", v)
				}
				if gw.To4() != nil {
					gw = gw.To4()
				}
			case "nic":
				i, err := strconv.ParseInt(v, 10, 32)
				if err != nil {
					fail("bad nic %q in route statement, must be int32", v)
				}
				nic = tcpip.NICID(i)
			default:
				fail("bad route keyword %q", k)
			}
		}
		if nic == 0 {
			fail("missing nic in route statement")
		}
		if ipnet == nil {
			fail("missing ipnet in route statement")
		}
		if gw == nil && !haveGW {
			fail("missing gw in route statement")
		}
		if _, ok := nics[nic]; !ok {
			fail("unknown nic %v, define nic before using", nic)
		}
		route := tcpip.Route{
			Destination: tcpip.Address(ipnet.IP),
			Mask:        tcpip.AddressMask(ipnet.Mask),
			Gateway:     tcpip.Address(gw),
			NIC:         nic,
		}
		if verbose {
			log.Printf("netstack: adding route id=%d ipnet=%s gw=%s", nic, ipnet, gw)
		}
		routes = append(routes, route)
	}

	parseDNS := func(args [][2]string) {
		for _, p := range args {
			k, v := p[0], p[1]
			switch k {
			case "ip":
				_, _, err := SplitHostPort(v)
				if err != nil {
					// Assuming the error is "missing port", the resolver expects a host:port.
					v = JoinHostPort(v, "53")
				}
				netstackNameservers = append(netstackNameservers, v)
			default:
				fail("bad dns keyword %q", k)
			}
		}
		if verbose {
			log.Printf("netstack: dns %v", netstackNameservers)
		}
	}

	parseLine := func(line string) {
		t := strings.Split(line, " ")
		args := parseArgs(t[1:])
		switch t[0] {
		case "verbose":
			verbose = true
		case "nic":
			parseNIC(args)
		case "ip":
			parseIP(args)
		case "route":
			parseRoute(args)
		case "dns":
			parseDNS(args)
		default:
			fail("unknown keyword %q", t[0])
		}
	}

	for _, line := range strings.Split(gonet, ";") {
		line = strings.TrimSpace(line)
		parseLine(line)
	}

	protocols := []string{}
	if haveIP4 {
		protocols = append(protocols, ipv4.ProtocolName)
	}
	if haveIP6 {
		protocols = append(protocols, ipv6.ProtocolName)
	}
	if haveEther {
		protocols = append(protocols, arp.ProtocolName)
	}
	s := stack.New(protocols, []string{tcp.ProtocolName, udp.ProtocolName}, stack.Options{})

	for nic, link := range nics {
		if err := s.CreateNIC(nic, link); err != nil {
			fail("creating nic: %v", err)
		}
	}

	for _, a := range ipaddrs {
		num := ipv4.ProtocolNumber
		if len(a.ip) == 16 {
			num = ipv6.ProtocolNumber
		}
		if err := s.AddAddress(a.nic, num, tcpip.Address(a.ip)); err != nil {
			fail("adding address %s to nic %d in stack: %v", a.ip, a.nic, err)
		}
	}

	for nic, _ := range ethers {
		if err := s.AddAddress(nic, arp.ProtocolNumber, "arp"); err != nil {
			fail("adding arp to stack: %v", err)
		}
	}

	s.SetRouteTable(routes)

	netstack = s
}
