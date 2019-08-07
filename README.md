# Modifications

This branch has an experiment for having a runtime without access
to files. The changes only apply to openbsd/amd64, to be build on
any other arch (self-hosting won't work without files).

## Changes

- syscall.Open returns ENOTSUP
- os.Getwd always returns "/"
- os.Open always returns ENOTSUP
- use getentropy() syscall during runtime init, not openening /dev/urandom. crypto/rand already uses getentropy.

- More disabled: other file-related syscalls (stat, umask, etc), fork & exec, ioctl, bpf (was deprecated)
- Added builtin files, registered with os.AddFile. can be used for reading /etc/resolv.conf, or /etc/ssl/cert.pem. Also added os.PrintOpen(bool), to toggle printing Open's, useful while developing.
- time.SetTimezoneDB lets you set contents of the lib/time/zoneinfo.zip. For now, timezone config can be done through TZ, eg TZ=Europe/Amsterdam. Should perhaps just add always include the zip file in the time package.

## Notes
- Many changes are for all of openbsd. Have to revisit later, but it should get a separate architecture.
- Many syscall numbers have been removed. Mostly to catch uses of them. Might want to add them back later.
- Should reenable API checking, with the correct baseline/exceptions.

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
