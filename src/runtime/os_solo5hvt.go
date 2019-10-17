// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"unsafe"
)

type mOS struct{}
type sigset uint32
type gsignalStack struct{}

func osinit() {
	ncpu = 1
	physPageSize = 4 * 1024
}

//go:nosplit
func getRandomData(r []byte) {
	// xxx todo bug, need to actually set random data!
}

const _NSIG = 0

func signame(sig uint32) string {
	return ""
}

func initsig(preinit bool) {
}

func goenvs() {
	solo5envs()
}

// Called to initialize a new m (including the bootstrap m).
// Called on the parent thread (main thread in case of bootstrap), can allocate memory.
func mpreinit(mp *m) {
	mp.gsignal = malg(2 * 1024)
	mp.gsignal.m = mp
}

//go:nosplit
func msigsave(mp *m) {
}

//go:nosplit
func msigrestore(sigmask sigset) {
}

//go:nosplit
//go:nowritebarrierrec
func clearSignalHandlers() {
}

//go:nosplit
func sigblock() {
}

// Called to initialize a new m (including the bootstrap m).
// Called on the new thread, can not allocate memory.
func minit() {
}

// Called from dropm to undo the effect of an minit.
//go:nosplit
func unminit() {
}

func crash() {
	*(*int32)(nil) = 0
}

func setProcessCPUProfiler(hz int32) {}
func setThreadCPUProfiler(hz int32)  {}
func sigdisable(uint32)              {}
func sigenable(uint32)               {}
func sigignore(uint32)               {}

//go:linkname os_sigpipe os.sigpipe
func os_sigpipe() {
	throw("too many writes on closed pipe")
}

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
	MemSize      uintptr // memory size in bytes
	KernelEnd    uintptr // address of the end of kernel
	CpuCycleFreq uint64  // CPU cycle counter frequency, Hz
	Cmdline      uintptr // command-line, c string
	Manifest     uintptr // address of application manifest
}

var solo5BootInfo *bootInfo
var Solo5BootInfo *bootInfo

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

// xxx do not export
//go:nosplit
func Solo5Write(buf []byte) {
	var arg = struct {
		data   uintptr
		length uint64
	}{uintptr(unsafe.Pointer(&buf[0])), uint64(len(buf))}
	outl(hypercallPuts, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(&arg)
}

//go:nosplit
func solo5Putp(p uintptr, n int) {
	var arg = struct {
		data   uintptr
		length uint64
	}{p, uint64(n)}
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
		ret      int64
	}{nsec, 0, 0}
	outl(hypercallPoll, uintptr(unsafe.Pointer(&arg)))
	return arg.readySet, arg.ret
}

//go:nosplit
func solo5Blkwrite(handle, offset uint64, data []byte) int64 {
	var arg = struct {
		// in
		handle uint64
		offset uint64
		data   uintptr
		length int64

		// out
		ret int64
	}{handle, offset, uintptr(unsafe.Pointer(&data[0])), int64(len(data)), -1}
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
		data   uintptr

		// in/out
		length int64

		// out
		ret int64
	}{handle, offset, uintptr(unsafe.Pointer(&data[0])), int64(len(data)), 0}
	outl(hypercallBlkread, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(data)
	return arg.length, arg.ret
}

//go:nosplit
func solo5Netwrite(handle uint64, data []byte) int64 {
	var arg = struct {
		// in
		handle uint64
		data   uintptr
		length int64

		// out
		ret int64
	}{handle, uintptr(unsafe.Pointer(&data[0])), int64(len(data)), -1}
	outl(hypercallNetwrite, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(data)
	return arg.ret
}

//go:nosplit
func solo5Netread(handle uint64, data []byte) (int64, int64) {
	var arg = struct {
		// in
		handle uint64
		data   uintptr

		// in/out
		length int64

		// out
		ret int64
	}{handle, uintptr(unsafe.Pointer(&data[0])), int64(len(data)), 0}
	outl(hypercallNetread, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(data)
	return arg.length, arg.ret
}

//go:nosplit
func exit(code int32) {
	var arg = struct {
		// in
		cookie     uintptr
		exitStatus int64
	}{0, int64(code)}
	outl(hypercallHalt, uintptr(unsafe.Pointer(&arg)))
	KeepAlive(&arg)
}

var buffer [512]byte
var fmtbuf [32]byte

var (
	tscBase        uint64
	nanotimeBase   uint64
	walltimeOffset uint64 // offset from nanotime

	// Nanotime is calculated as (tsc * tscMult) >> tscShift
	// With "tsc" being the current tsc.
	// And tscMult & tscShift calculated during initialization.
	tscMult  uint32
	tscShift uint8
)

func nanotime() int64 {
	tsc := uint64(cputicks())
	tscDelta := tsc - tscBase
	nanotimeBase += (tscDelta * uint64(tscMult)) >> tscShift
	tscBase = tsc
	return int64(nanotimeBase)
}

func walltime() (sec int64, nsec int32) {
	ns := uint64(nanotime()) + walltimeOffset
	sec = int64(ns / 1e9)
	nsec = int32(ns % 1e9)
	return
}

//go:nosplit
func solo5init(bi *bootInfo) {
	solo5BootInfo = bi
	Solo5BootInfo = bi
	memoryNext = bi.KernelEnd
	memoryEnd = bi.MemSize

	// Initialize time using TSC, from solo5.
	tscShift = 32
	const nanoseconds = 1e9
	for tscShift > 0 && tscMult == 0 {
		tmp := (nanoseconds << tscShift) / bi.CpuCycleFreq
		if (tmp & 0xffffffff00000000) == 0 {
			tscMult = uint32(tmp)
		} else {
			tscShift--
		}
	}
	if tscMult == 0 {
		solo5Puts("bad CpuCycleFreq\n")
		exit(1)
	}

	tscBase = uint64(cputicks())
	nanotimeBase = (tscBase * uint64(tscMult)) >> tscShift // todo: use something like the mul64_32() from solo5?
	walltimeOffset = solo5Walltime() - nanotimeBase

	solo5Puts("solo5 init...\n")

	/*
		_, ret := solo5Poll(10*1000*1000)
		if ret == 0 {
			solo5Puts("ret 0\n")
		} else {
			solo5Puts("ret >0\n")
		}

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

		solo5Puts("time: ")
		nsec := solo5Walltime()
		Write(decimal(nsec))

		solo5Puts("bootInfo:\n")
		solo5Puts("memSize: ")
		Write(hex(bi.MemSize))
		solo5Puts("kernelEnd: ")
		Write(hex(uint64(bi.KernelEnd)))
		solo5Puts("cpuCycleFreq: ")
		Write(hex(bi.CpuCycleFreq))

		slen := 0
		for *(*byte)(unsafe.Pointer(bi.Cmdline + uintptr(slen))) != 0 {
			slen++
		}
		solo5Puts("cmdline: ")
		solo5Putp(bi.Cmdline, slen)
		solo5Puts("\n")

		solo5Puts("manifest: ")
		Write(hex(uint64(bi.Manifest)))

		// handle, ok := solo5Lookup("blk0")
		manifest := (*manifest)(unsafe.Pointer(bi.Manifest))
		Write(decimal(uint64(manifest.version)))
		Write(decimal(uint64(manifest.nentries)))
		Write(hex(uint64(uintptr((unsafe.Pointer(&manifest.entries))))))
		for i := uint32(1); i < manifest.nentries; i++ {
			p := uintptr(unsafe.Pointer(&manifest.entries))
			p += uintptr(i * mftEntrySize)
			e := (*mftEntry)(unsafe.Pointer(p))
			solo5Puts("device: ")
			Write(decimal(uint64(e.etype)))
			solo5Putp(uintptr(unsafe.Pointer(&e.name)), 68)
			solo5Puts("\n")
		}

		solo5Halt(nil, 3)
	*/
}

func solo5LookupBlock(bi *bootInfo, name string) (handle uint64, info *mftBlockBasic) {
	manifest := (*manifest)(unsafe.Pointer(bi.Manifest))
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
	manifestDevNetBasic   = 2
	manifestReservedFrist = 1 << 30
)

type manifest struct {
	version  uint32
	nentries uint32
	entries  uint64 // mftEntry
}

type mftEntry struct {
	name  [68]byte // c string
	etype uint32
	info  uintptr // either mftBlockBasic or mftNetBasic
}

type mftBlockBasic struct {
	capacity  uint64
	blockSize uint16
}

type mftNetBasic struct {
	mac [6]byte
	mtu uint16
}

// xxx
var noenv = []string{}
var noargs = []string{}

func split(s string) (t []string) {
	for s != "" {
		i := 0
		for i < len(s) && s[i] != ' ' {
			i++
		}
		t = append(t, s[:i])
		if i < len(s) {
			i++
		}
		s = s[i:]
	}
	return
}

func solo5envs() {
	// solo5BootInfo.Cmdline looks like: -env KEY=value -env KEY=value -- arg1 arg2 ...

	cmdline := gostring((*byte)(unsafe.Pointer(solo5BootInfo.Cmdline)))

	t := split(cmdline)
	envs = []string{}
	for len(t) > 0 && t[0] == "-env" {
		if len(t) == 1 {
			panic("bad cmdline")
		}
		envs = append(envs, t[1])
		t = t[2:]
	}
	if len(t) > 0 && t[0] == "--" {
		t = t[1:]
	}
	argslice = append([]string{"solo5hvt"}, t...)
}

/*
//go:nosplit
func readblock() {
	block := buffer[:]
	n, ret := solo5Blkread(1, 0, block)
	if ret != 0 {
		solo5Puts("read failed\n")
		block[0] = '0' + byte(ret)
		block[1] = '\n'
		Write(block[:2])
	} else {
		solo5Puts("read:\n")
		Write(block[:n])
		solo5Puts("\n")
	}
}
*/

// replacements for stubs2.go
func read(fd int32, p unsafe.Pointer, n int32) int32 {
	throw("read")
	return 0
}

func closefd(fd int32) int32 {
	throw("closefd")
	return 0
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
