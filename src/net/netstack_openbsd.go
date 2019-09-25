package net

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/header"
	"github.com/google/netstack/tcpip/link/fdbased"
	"github.com/google/netstack/tcpip/link/sniffer"
	"github.com/google/netstack/tcpip/network/arp"
	"github.com/google/netstack/tcpip/network/ipv4"
	"github.com/google/netstack/tcpip/network/ipv6"
	"github.com/google/netstack/tcpip/stack"
	"github.com/google/netstack/tcpip/transport/icmp"
	"github.com/google/netstack/tcpip/transport/tcp"
	"github.com/google/netstack/tcpip/transport/udp"
)

var (
	netstack            *stack.Stack
	netstackNameservers []string
	errStack            = errors.New("no network stack")
)

func splitPair(s string) (string, string) {
	for i, c := range s{
		if c == '=' {
			return s[:i], s[i+1:]
		}
	}
	return s, ""
}

func init() {
	gonet := os.Getenv("GONET")
	if gonet == "" {
		// no stack configured
		return
	}

	rand.Seed(time.Now().UnixNano())

	fail := func(args ...string) {
		s := ""
		for _, v := range args {
			s += v
		}
		panic(errors.New(s))
	}

	// Example of what we are trying to parse:
	//
	//	verbose
	// 	nic id=1 ether=12:12:12:23:23:23 mtu=1500 fd=3 dev=/dev/tap0 sniff=true
	// 	ip nic=1 addr=1.2.3.4/24 addr=2.3.4.5/32 addr=ff:aa.../64
	// 	route nic=1 ipnet=ip/mask
	// 	route nic=1 ipnet=ip/mask gw=ip...
	// 	dns ip=1.2.3.4 ip=8.8.8.8

	parseArgs := func(pairs []string) [][2]string {
		args := make([][2]string, len(pairs))
		for i, p := range pairs {
			k, v := splitPair(p)
			args[i] = [...]string{k, v}
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
		ethers                      = map[tcpip.NICID]HardwareAddr{}
		ipaddrs                     []nicip
		routes                      []tcpip.Route
		haveIP4, haveIP6, haveEther bool
	)

	xparseNIC := func(s string) tcpip.NICID {
		v, _, ok := dtoi(s)
		if !ok {
			fail("invalid nic id ", s)
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
					fail("parsing ether ", v, ": ", err.Error())
				}
				if len(mac) != 6 {
					fail("invalid ether length, must be six bytes, saw ", v, " of ", itoa(len(mac)))
				}
				ether = mac
				haveEther = true
			case "mtu":
				var ok bool
				mtu, _, ok = dtoi(v)
				if !ok {
					fail("parsing mtu ", v)
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
					fail("bad sniff ", v)
				}
			default:
				fail("bad route keyword ", k)
			}
		}
		if nic == 0 {
			fail("nic: missing id")
		}
		if _, ok := nics[nic]; ok {
			fail("nic: duplicate nic ", itoa(int(nic)))
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
			fail("making new fd-based link: ", err.Error())
		}
		if sniff {
			linkID = sniffer.New(linkID)
		}
		nics[nic] = linkID
		if ether != nil {
			ethers[nic] = ether
		}
		if verbose {
			print("netstack: adding nic id=", nic, " ether=", ether.String(), " mtu=", mtu, " fd=", fd, " dev=", dev, " sniff=", sniff, "\n")
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
					fail("bad addr value ", v)
				}
				ips = append(ips, ip)
				if ip.To4() != nil {
					haveIP4 = true
				} else {
					haveIP6 = true
				}
			default:
				fail("bad ip keyword ", k)
			}
		}
		if nic == 0 {
			fail("missing nic in ip statement")
		}
		if _, ok := nics[nic]; !ok {
			fail("unknown nic ", itoa(int(nic)), " in ip statement, define it before use")
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
				print("netstack: adding ip nic=", nic, " addr=", a.ip.String(), "\n")
			}
			ipaddrs = append(ipaddrs, a)
		}
	}

	parseRoute := func(args [][2]string) {
		var (
			nic   tcpip.NICID
			ipnet *IPNet
			gw    IP
		)
		for _, p := range args {
			k, v := p[0], p[1]
			switch k {
			case "ipnet":
				var err error
				_, ipnet, err = ParseCIDR(v)
				if err != nil {
					fail("bad ipnet ", v, " in route statement: ", err.Error())
				}
			case "gw":
				if v == "" {
					gw = nil
					continue
				}
				gw = ParseIP(v)
				if gw == nil {
					fail("bad gw ", v, " in route statement")
				}
				gw4 := gw.To4()
				if gw4 != nil {
					gw = gw4
				}
			case "nic":
				i, _, ok := dtoi(v)
				if !ok {
					fail("bad nic ", v, " in route statement, must be int")
				}
				nic = tcpip.NICID(i)
			default:
				fail("bad route keyword ", k)
			}
		}
		if nic == 0 {
			fail("missing nic in route statement")
		}
		if ipnet == nil {
			fail("missing ipnet in route statement")
		}
		if _, ok := nics[nic]; !ok {
			fail("unknown nic ", itoa(int(nic)), " define nic before using")
		}
		subnet, err := tcpip.NewSubnet(tcpip.Address(ipnet.IP), tcpip.AddressMask(ipnet.Mask))
		if err != nil {
			fail("bad subnet ", ipnet.String(), ": ", err.Error())
		}
		route := tcpip.Route{
			Destination: subnet,
			Gateway:     tcpip.Address(gw),
			NIC:         nic,
		}
		if verbose {
			sgw := ""
			if gw != nil {
				sgw = gw.String()
			}
			print("netstack: adding route id=", nic, " ipnet=", ipnet.String(), " gw=", sgw, "\n")
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
				fail("bad dns keyword ", k)
			}
		}
		if verbose {
			s := "netstack: dns"
			for _, ns := range netstackNameservers {
				s += " " + ns
			}
			s += "\n"
			print(s)
		}
	}

	parseLine := func(line string) {
		t := splitAtBytes(line, " ")
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
			fail("unknown keyword ", t[0])
		}
	}

	for _, line := range splitAtBytes(gonet, ";") {
		line = string(trimSpace([]byte(line)))
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
		if haveIP4 {
			protocols = append(protocols, arp.ProtocolName)
		}
	}
	transports := []string{tcp.ProtocolName, udp.ProtocolName}
	if haveIP4 {
		transports = append(transports, icmp.ProtocolName4)
	}
	if haveIP6 {
		transports = append(transports, icmp.ProtocolName6)
	}
	s := stack.New(protocols, transports, stack.Options{})

	for nic, link := range nics {
		if err := s.CreateNIC(nic, link); err != nil {
			fail("creating nic: ", err.String())
		}
	}

	solDone := map[tcpip.NICID]struct{}{}
	for _, a := range ipaddrs {
		num := ipv4.ProtocolNumber
		if len(a.ip) == 16 {
			num = ipv6.ProtocolNumber
		}
		if err := s.AddAddress(a.nic, num, tcpip.Address(a.ip)); err != nil {
			fail("adding address ", a.ip.String(), " to nic ", itoa(int(a.nic)), " in stack: ", err.String())
		}

		if len(a.ip) == 16 {
			snaddr := header.SolicitedNodeAddr(tcpip.Address(a.ip))
			if verbose {
				print("netstack: adding ndp address", snaddr.String(), "\n")
			}
			s.AddAddress(a.nic, num, snaddr)
			if _, ok := solDone[a.nic]; !ok {
				ether, ok := ethers[a.nic]
				if ok {
					lladdr := header.LinkLocalAddr(tcpip.LinkAddress(ether))
					snlladdr := header.SolicitedNodeAddr(lladdr)
					/*
						_, localnet, _ := ParseCIDR("fe80::/10")
						subnet, err := tcpip.NewSubnet(tcpip.Address(localnet.IP), tcpip.AddressMask(localnet.Mask))
						if err != nil {
							fail("making local subnet: ", err.Error())
						}
						route := tcpip.Route{
							Destination: subnet,
							NIC: a.nic,
						}
						routes = append(routes, route)
					*/
					if verbose {
						print("netstack: adding link local addr ", lladdr.String(), " for nic ", a.nic, "\n")
						print("netstack: adding link local ndp address ", snlladdr.String(), " for nic ", a.nic, "\n")
					}
					s.AddAddress(a.nic, num, lladdr)
					s.AddAddress(a.nic, num, snlladdr)
				}
				solDone[a.nic] = struct{}{}
			}
		}
	}

	for nic, _ := range ethers {
		if err := s.AddAddress(nic, arp.ProtocolNumber, "arp"); err != nil {
			fail("adding arp to stack: ", err.String())
		}
	}

	s.SetRouteTable(routes)

	netstack = s
}
