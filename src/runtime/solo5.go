// +build openbsd

package runtime

import (
	"unsafe"
)

func outl(dx uint32, ax uintptr)

const (
	hypercallPioBase = 0x500 + iota
	hypercallWalltime
	hypercallPuts
	hypercallPoll
	hypercallBlkwrite
	hypercallBlkread
	hypercallNetwrite
	hypercallNetread
	hypercallHalt
)

const (
	s5ok = iota
	s5again
	s5invalid
	s5unspec
)

type bootInfo struct {
	memSize uint64 // memory size in bytes
	kernelEnd uintptr // address of the end of kernel
	cpuCycleFreq uint64 // CPU cycle counter frequency, Hz
	cmdline uintptr	// command-line, c string
	manifest uintptr	// address of application manifest
}

var solo5BootInfo *bootInfo

//go:nosplit
func solo5Walltime() (nsecs uint64) {
	outl(hypercallWalltime, uintptr(unsafe.Pointer(&nsecs)))
	return
}

//go:nosplit
func solo5Puts(s string) {
	outl(hypercallPuts, uintptr(unsafe.Pointer(&s)))
	KeepAlive(s)
}

//go:nosplit
func solo5Put(buf []byte) {
	var arg = struct {
		data uintptr
		length	uint64
	} { uintptr(unsafe.Pointer(&buf[0])), uint64(len(buf)) }
	outl(hypercallPuts, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(&arg)
}

//go:nosplit
func solo5Putp(p uintptr, n int) {
	var arg = struct {
		data uintptr
		length	uint64
	} { p, uint64(n) }
	outl(hypercallPuts, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(&arg)
}

//go:nosplit
func solo5Poll(nsec uint64) (uint64, int64) {
	var arg = struct {
		// in
		timeoutNsecs uint64

		// out
		readySet uint64
		ret int64
	} { nsec, 0, 0 }
	outl(hypercallPoll, uintptr(unsafe.Pointer(&arg)))
	return arg.readySet, arg.ret
}

//go:nosplit
func solo5Blkwrite(handle, offset uint64, data []byte) int64 {
	var arg = struct {
		// in
		handle uint64
		offset uint64
		data uintptr
		length int64

		// out
		ret int64
	} {handle, offset, uintptr(unsafe.Pointer(&data[0])), int64(len(data)), -1}
	outl(hypercallBlkwrite, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(data)
	return arg.ret
}

//go:nosplit
func solo5Blkread(handle, offset uint64, data []byte) (int64, int64) {
	var arg = struct {
		// in
		handle uint64
		offset uint64
		data uintptr

		// in/out
		length int64

		// out
		ret int64
	} { handle, offset, uintptr(unsafe.Pointer(&data[0])), int64(len(data)), 0 }
	outl(hypercallBlkread, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(data)
	return arg.length, arg.ret
}


//go:nosplit
func solo5Netwrite(handle uint64, data []byte) int64 {
	var arg = struct {
		// in
		handle uint64
		data uintptr
		length int64

		// out
		ret int64
	} {handle, uintptr(unsafe.Pointer(&data[0])), int64(len(data)), -1}
	outl(hypercallNetwrite, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(data)
	return arg.ret
}

//go:nosplit
func solo5Netread(handle uint64, data []byte) (int64, int64) {
	var arg = struct {
		// in
		handle uint64
		data uintptr

		// in/out
		length int64

		// out
		ret int64
	} { handle, uintptr(unsafe.Pointer(&data[0])), int64(len(data)), 0 }
	outl(hypercallNetread, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(data)
	return arg.length, arg.ret
}

//go:nosplit
func solo5Halt(cookie []byte, status int64) {
	var arg = struct {
		// in
		cookie uintptr
		exitStatus int64
	} { 0, status }
	if len(cookie) > 0 {
		arg.cookie = uintptr(unsafe.Pointer(&cookie[0]))
	}
	outl(hypercallHalt, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(&arg)
}

var buffer [512]byte
var fmtbuf [32]byte

//go:nosplit
func solo5init(bi *bootInfo) {
	solo5BootInfo = bi
	memoryEnd = bi.kernelEnd

	solo5Puts("solo5 init...\n")

/*
	_, ret := solo5Poll(10*1000*1000)
	if ret == 0 {
		solo5Puts("ret 0\n")
	} else {
		solo5Puts("ret >0\n")
	}
*/

	readblock()

	decimal := func(v uint64) []byte {
		e := len(fmtbuf)
		e--
		fmtbuf[e] = '\n'
		for v > 0 {
			e--
			fmtbuf[e] = '0' + byte(v % 10)
			v /= 10
		}
		return fmtbuf[e:]
	}

	hex := func(v uint64) []byte {
		e := len(fmtbuf)
		e--
		fmtbuf[e] = '\n'
		for v > 0 {
			e--
			n := byte(v & 0xf)
			if n <= 9 {
				fmtbuf[e] = '0' + n
			} else {
				fmtbuf[e] = 'a' + n - 10
			}
			v >>= 4
		}
		return fmtbuf[e:]
	}

/*
	solo5Puts("time: ")
	nsec := solo5Walltime()
	solo5Put(decimal(nsec))

	solo5Puts("bootInfo:\n")
	solo5Puts("memSize: ")
	solo5Put(hex(bi.memSize))
	solo5Puts("kernelEnd: ")
	solo5Put(hex(uint64(bi.kernelEnd)))
	solo5Puts("cpuCycleFreq: ")
	solo5Put(hex(bi.cpuCycleFreq))

	slen := 0
	for *(*byte)(unsafe.Pointer(bi.cmdline + uintptr(slen))) != 0 {
		slen++
	}
	solo5Puts("cmdline: ")
	solo5Putp(bi.cmdline, slen)
	solo5Puts("\n")
*/

	solo5Puts("manifest: ")
	solo5Put(hex(uint64(bi.manifest)))

	// handle, ok := solo5Lookup("blk0")
	manifest := (*manifest)(unsafe.Pointer(bi.manifest))
	solo5Put(decimal(uint64(manifest.version)))
	solo5Put(decimal(uint64(manifest.nentries)))
	solo5Put(hex(uint64(uintptr((unsafe.Pointer(&manifest.entries))))))
	for i := uint32(1); i < manifest.nentries; i++ {
		p := uintptr(unsafe.Pointer(&manifest.entries))
		p += uintptr(i * mftEntrySize)
		e := (*mftEntry)(unsafe.Pointer(p))
		solo5Puts("device: ")
		solo5Put(decimal(uint64(e.etype)))
		solo5Putp(uintptr(unsafe.Pointer(&e.name)), 68)
		solo5Puts("\n")
	}

	// solo5Halt(nil, 3)
}

func solo5LookupBlock(bi *bootInfo, name string) (handle uint64, info *mftBlockBasic) {
	manifest := (*manifest)(unsafe.Pointer(bi.manifest))
	for i := uint32(1); i < manifest.nentries; i++ {
		p := uintptr(unsafe.Pointer(&manifest.entries))
		p += uintptr(i * mftEntrySize)
		e := (*mftEntry)(unsafe.Pointer(p))
		if e.etype != manifestDevBlockBasic {
			continue
		}
		// xxx
	}
	return 0, nil
}

const mftEntrySize = 68 + 4 + 16 + 8 + 1 + 7

const (
	manifestDevBlockBasic = 1
	manifestDevNetBasic = 2
	manifestReservedFrist = 1<<30
)

type manifest struct {
	version uint32
	nentries uint32
	entries uint64 // mftEntry
}

type mftEntry struct {
	name [68]byte	// c string
	etype uint32
	info uintptr	// either mftBlockBasic or mftNetBasic
}

type mftBlockBasic struct {
	capacity uint64
	blockSize uint16
}

type mftNetBasic struct {
	mac [6]byte
	mtu uint16
}

// xxx
var noenv = []string{}
var noargs = []string{}

func solo5envs() {
	// xxx parse solo5BootInfo.cmdline
	envs = noenv
	argslice = noargs
}

//go:nosplit
func readblock() {
	block := buffer[:]
	n, ret := solo5Blkread(1, 0, block)
	if ret != 0 {
		solo5Puts("read failed\n")
		block[0] = '0' + byte(ret)
		block[1] = '\n'
		solo5Put(block[:2])
	} else {
		solo5Puts("read:\n")
		solo5Put(block[:n])
		solo5Puts("\n")
	}
}

// replacements for stubs2.go
func read(fd int32, p unsafe.Pointer, n int32) int32 {
	throw("read")
	return 0
}

func closefd(fd int32) int32 {
	throw("closefd")
	return 0
}

func exit(code int32) {
	solo5Halt(nil, int64(code))
}

func usleep(usec uint32) {
}

func write(fd uintptr, p unsafe.Pointer, n int32) int32 {
	solo5Putp(uintptr(p), int(n))
	KeepAlive(p)
	return 0
}

func open(name *byte, mode, perm int32) int32 {
	throw("open")
	return 0
}

// return value is only set on linux to be used in osinit()
func madvise(addr unsafe.Pointer, n uintptr, flags int32) int32 {
	throw("madvise")
	return 0
}

// exitThread terminates the current thread, writing *wait = 0 when
// the stack is safe to reclaim.
//
func exitThread(wait *uint32) {
	throw("exitThread")
}

func osyield() {
	throw("osyield")
}

func newosproc(mp *m) {
	throw("newosproc")
}
