# vmgo

NOTE: this is work in progress.

The goal of vmgo is to get to a go toolchain that can compile
existing Go code (with no or otherwise minimal changes) to a
standalone virtual image that can run in a (minimal) virtual machine
monitor.

As an example, something like this should be possible:

	cd some/existing/project
	GOOS=somevm GOARCH=amd64 go build -o project.img
	some-minimal-vmm -net xxxnetconfig project.img

This branch started as an experiment for having a runtime without
access to files. The changes apply to openbsd/amd64, with
cross-compilation (self-hosting won't work without files).

Lots of functionality has to be stripped down, or replaced with a
builtin implementation, because there is no OS to provide the
functionality. Files. Network stack. Processes, signals. Most
remaining system calls.

The initial target is solo5.

## changes

(this list is probably incomplete)

	- many "syscall" functions return ENOTSUP
	- file system-related typically return errors. os.Open() only works on files that were added with the new os.AddFile(path, data), for adding /etc/resolv.conf, /etc/ssl/cert.pem, etc.  os.PrintOpen(bool) toggles printing opens, for debugging.
	- getentropy() is used during runtime init, not /dev/urandom (crypto/rand already uses getentropy)
	- no cgo, probably necessary anyway, but also gets rid of one more variable.
	- added time.SetTimezoneDB to set fake contents for GOROOT+lib/time/zoneinfo.zip. For now, timezone config can be done through TZ, eg TZ=Europe/Amsterdam. Should perhaps just add always include the zip file in the time package.

## file access

The runtime does surprisingly few file opens. These might be opened on openbsd, and replacements need to be provided.

	- /etc/ssl/cert.pem, src/crypto/x509/root_bsd.go
	- /etc/mime.types, src/mime/type_unix.go; /usr/share/misc/mime.types on openbsd
	- /etc/hosts, in net for resolving, in src/net/hook.go
	- /etc/services, in readServices, src/net/port_unix.go
	- /etc/protocols, in readProtocols, src/net/lookup_unix.go
	- /etc/user, /etc/group, in lookupGroup, lookupUser, in src/os/user/lookup_unix.go
	- /dev/log, /var/run/syslog, /var/run/log in src/log/syslog/syslog_unix.go, net.Dial with unix or unixgram as parameter.

## todo

	- add a "net" backend that talks to a tun device (as a file descriptor for now).
	- remove process creation from the runtime for the openbsd arch. take a hint from wasm, which is single process.
	- some way to pass & parse variables when there are no more regular env vars.
	- remove more system calls. see how setting TLS for "g" can be stripped.
	- adjust to target solo5 (or something else by that time).

	- restore API check during toolchain build.
	- don't target "openbsd", but add a new arch.


(original Go README below)


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
