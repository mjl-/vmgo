package fdbased

import (
	"syscall"
	"bufio"
	"fmt"
	"os"
	"unsafe"
)


var dumpEnabled = os.Getenv("NETSTACKDUMP") != ""

func dumpv(what string, iovecs []syscall.Iovec, n int) {
	if !dumpEnabled {
		return
	}
	b := bufio.NewWriter(os.Stderr)
	fmt.Fprintf(b, "%s packet:\n", what)
	iov := make([]syscall.Iovec, len(iovecs))
	copy(iov, iovecs)
	for n > 0 {
		for i := 0; i < 16 && n > 0; i++ {
			n--
			v := *iov[0].Base
			iov[0].Len--
			if iov[0].Len == 0 {
				iov = iov[1:]
			} else {
				iov[0].Base = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(iov[0].Base))+1))
			}
			fmt.Fprintf(b, "%02x", v)
			if i % 2 == 1 {
				fmt.Fprint(b, " ")
			}
		}
		fmt.Fprint(b, "\n")
	}
	fmt.Fprint(b, "\n")
	b.Flush()
}

func dump(what string, buf []byte) {
	if !dumpEnabled {
		return
	}
	b := bufio.NewWriter(os.Stderr)
	fmt.Fprintf(b, "%s packet:\n", what)
	for len(buf) > 0 {
		for i := 0; i < 8 && len(buf) > 0; i++ {
			n := 2
			if len(buf) == 1 {
				n = 1
			}
			fmt.Fprintf(b, "%x ", buf[:n])
			buf = buf[n:]
		}
		fmt.Fprint(b, "\n")
	}
	fmt.Fprint(b, "\n")
	b.Flush()
}
