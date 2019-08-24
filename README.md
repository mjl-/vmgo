# "net" with github.com/google/netstack

First build this toolchain and make sure you're using the new toolchain.

Build your program with the "netstack" build tag:

	go build -tags netstack webapp.go

Set up your networking. See the example below for configuring a bridge and a tap device "tap0".

Now run your program with a GONET environment variable:

	GONET='verbose; nic id=1 ether=fe:e1:ba:d0:33:33 mtu=1500 dev=tap0 sniff=true; ip nic=1 addr=192.168.1.100/24; route nic=1 ipnet=0.0.0.0/0 gw=192.168.1.1; dns ip=8.8.8.8' ./webapp

Your program is now running with a pure Go built-in networking stack. You will see a trace of all packets.

The GONET variable configures the network, with semicolon-separated statements. "verbose" causes the config to be printed at startup. "nic" configures the NIC, with "sniff=true" enabling the netstack sniffer to print summaries of all incoming & outgoing packets. "ip" adds ip addresses. "route" adds routes.


## Example networking setup

	ip link add bridge0 type bridge
	ip link set dev bridge0 up

	ip addr flush dev enp0s31f6  # fill in your ethernet device. this command removes the existing ip's.
	ip link set enp0s31f6 master bridge0  # add ethernet to bridge
	ip addr add 192.168.1.200/24 dev bridge0  # fill in your local ip address, use a different address in GONET; optional
	ip route add default via 192.168.1.1  # fill in your local gateway; optional

	ip tuntap add dev tap0 mode tap user mjl  # use your username instead of "mjl"
	ip link set dev tap0 address fe:e1:ba:d0:11:11  # configure fake mac address. use a different fake mac address in your GONET
	ip link set tap0 master bridge0
	ip link set dev tap0 up

	sysctl net.ipv4.ip_forward=1  # only needed when using "tun" device, not needed for "tap"


# The Go Programming Language

Go is an open source programming language that makes it easy to build simple,
reliable, and efficient software.

![Gopher image](doc/gopher/fiveyears.jpg)
*Gopher image by [Renee French][rf], licensed under [Creative Commons 3.0 Attributions license][cc3-by].*

Our canonical Git repository is located at https://go.googlesource.com/go.
There is a mirror of the repository at https://github.com/golang/go.

Unless otherwise noted, the Go source files are distributed under the
BSD-style license found in the LICENSE file.

### Download and Install

#### Binary Distributions

Official binary distributions are available at https://golang.org/dl/.

After downloading a binary release, visit https://golang.org/doc/install
or load [doc/install.html](./doc/install.html) in your web browser for installation
instructions.

#### Install From Source

If a binary distribution is not available for your combination of
operating system and architecture, visit
https://golang.org/doc/install/source or load [doc/install-source.html](./doc/install-source.html)
in your web browser for source installation instructions.

### Contributing

Go is the work of thousands of contributors. We appreciate your help!

To contribute, please read the contribution guidelines:
	https://golang.org/doc/contribute.html

Note that the Go project uses the issue tracker for bug reports and
proposals only. See https://golang.org/wiki/Questions for a list of
places to ask questions about the Go language.

[rf]: https://reneefrench.blogspot.com/
[cc3-by]: https://creativecommons.org/licenses/by/3.0/
