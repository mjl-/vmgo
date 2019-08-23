package net

import (
	"errors"
	"fmt"
)

// If the ifindex is zero, interfaceTable returns mappings of all
// network interfaces. Otherwise it returns a mapping of a specific
// interface.
func interfaceTable(ifindex int) ([]Interface, error) {
	if netstack == nil {
		return nil, errStack
	}
	var l []Interface

	for id, info := range netstack.NICInfo() {
		if ifindex == 0 || int(id) == ifindex {
			var flags Flags
			if info.Flags.Up {
				flags |= FlagUp
			}
			if info.Flags.Loopback {
				flags |= FlagLoopback
			}
			name := info.Name
			if name == "" {
				name = fmt.Sprintf("nic%d", id)
			}
			iface := Interface{
				Index:        int(id),
				MTU:          int(info.MTU),
				Name:         name,
				HardwareAddr: HardwareAddr(info.LinkAddress),
				Flags:        flags,
			}
			l = append(l, iface)
		}
	}
	if ifindex > 0 && len(l) == 0 {
		return nil, errors.New("interface not found")
	}
	return l, nil
}

// If the ifi is nil, interfaceAddrTable returns addresses for all
// network interfaces. Otherwise it returns addresses for a specific
// interface.
func interfaceAddrTable(ifi *Interface) ([]Addr, error) {
	for id, info := range netstack.NICInfo() {
		if int(id) == ifi.Index {
			var l []Addr
			for _, paddr := range info.ProtocolAddresses {
				// xxx netstack always has all bits set, because we add the address separately from the subnet.
				awp := paddr.AddressWithPrefix
				ip := IP(awp.Address)
				bits := 32
				if ip.To4() == nil {
					bits = 128
				}
				mask := CIDRMask(awp.PrefixLen, bits)
				addr := &IPNet{IP: ip, Mask: mask}
				l = append(l, addr)
			}
			return l, nil
		}
	}
	return nil, errors.New("interface not found")
}

// interfaceMulticastAddrTable returns addresses for a specific
// interface.
func interfaceMulticastAddrTable(ifi *Interface) ([]Addr, error) {
	return nil, notyet("interfaceMulticastAddrTable")
}
