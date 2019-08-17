// Copyright 2018 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fdbased

import (
	"syscall"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/buffer"
	"github.com/google/netstack/tcpip/header"
	"github.com/google/netstack/tcpip/link/rawfile"
	"github.com/google/netstack/tcpip/stack"
)

// BufConfig defines the shape of the vectorised view used to read packets from the NIC.
var BufConfig = []int{128, 256, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768}

// readVDispatcher uses readv() system call to read inbound packets and
// dispatches them.
type readVDispatcher struct {
	// fd is the file descriptor used to send and receive packets.
	fd int

	// e is the endpoint this dispatcher is attached to.
	e *endpoint

	// views are the actual buffers that hold the packet contents.
	views []buffer.View

	// iovecs are initialized with base pointers/len of the corresponding
	// entries in the views defined above, except when GSO is enabled then
	// the first iovec points to a buffer for the vnet header which is
	// stripped before the views are passed up the stack for further
	// processing.
	iovecs []syscall.Iovec
}

func newReadVDispatcher(fd int, e *endpoint) (linkDispatcher, error) {
	d := &readVDispatcher{fd: fd, e: e}
	d.views = make([]buffer.View, len(BufConfig))
	iovLen := len(BufConfig)
	if d.e.Capabilities()&stack.CapabilityGSO != 0 {
		iovLen++
	}
	iovLen++ // for openbsd tun leading 4 bytes
	d.iovecs = make([]syscall.Iovec, iovLen)
	return d, nil
}

func (d *readVDispatcher) allocateViews(bufConfig []int) {
	var vnetHdr [virtioNetHdrSize]byte
	vnetHdrOff := 0
	if d.e.Capabilities()&stack.CapabilityGSO != 0 {
		// The kernel adds virtioNetHdr before each packet, but
		// we don't use it, so so we allocate a buffer for it,
		// add it in iovecs but don't add it in a view.
		d.iovecs[0] = syscall.Iovec{
			Base: &vnetHdr[0],
			Len:  uint64(virtioNetHdrSize),
		}
		vnetHdrOff++
	}

	// openbsd adds 4 bytes with a "type" (AF_INET, AF_INET6) for tun files.
	// add an iovec (like the GSO-case above), but don't reference it otherwise.
	isTun := d.e.hdrSize == 0
	if isTun {
		var openbsdTunHdr [4]byte
		d.iovecs[vnetHdrOff] = syscall.Iovec{
			Base: &openbsdTunHdr[0],
			Len:  uint64(4),
		}
		vnetHdrOff++
	}

	for i := 0; i < len(bufConfig); i++ {
		if d.views[i] != nil {
			break
		}
		b := buffer.NewView(bufConfig[i])
		d.views[i] = b
		d.iovecs[i+vnetHdrOff] = syscall.Iovec{
			Base: &b[0],
			Len:  uint64(len(b)),
		}
	}
}

func (d *readVDispatcher) capViews(n int, buffers []int) int {
	c := 0
	for i, s := range buffers {
		c += s
		if c >= n {
			d.views[i].CapLength(s - (c - n))
			return i + 1
		}
	}
	return len(buffers)
}

// dispatch reads one packet from the file descriptor and dispatches it.
func (d *readVDispatcher) dispatch() (bool, *tcpip.Error) {
	d.allocateViews(BufConfig)

	n, err := rawfile.BlockingReadv(d.fd, d.iovecs)
	if err != nil {
		return false, err
	}
	dumpv("incoming", d.iovecs, n)
	if d.e.Capabilities()&stack.CapabilityGSO != 0 {
		// Skip virtioNetHdr which is added before each packet, it
		// isn't used and it isn't in a view.
		n -= virtioNetHdrSize
	}
	if n <= d.e.hdrSize {
		return false, nil
	}

	var (
		p             tcpip.NetworkProtocolNumber
		remote, local tcpip.LinkAddress
	)
	if d.e.hdrSize > 0 {
		eth := header.Ethernet(d.views[0])
		p = eth.Type()
		remote = eth.SourceAddress()
		local = eth.DestinationAddress()
	} else {
		// We don't get any indication of what the packet is, so try to guess
		// if it's an IPv4 or IPv6 packet.
		switch header.IPVersion(d.views[0]) {
		case header.IPv4Version:
			p = header.IPv4ProtocolNumber
		case header.IPv6Version:
			p = header.IPv6ProtocolNumber
		default:
			return true, nil
		}
	}

	used := d.capViews(n, BufConfig)
	vv := buffer.NewVectorisedView(n, d.views[:used])
	vv.TrimFront(d.e.hdrSize)

	d.e.dispatcher.DeliverNetworkPacket(d.e, remote, local, p, vv)

	// Prepare e.views for another packet: release used views.
	for i := 0; i < used; i++ {
		d.views[i] = nil
	}

	return true, nil
}
