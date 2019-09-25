package runtime

import (
	"unsafe"
)

// End of allocated memory. Next allocation comes after.
var memoryEnd uintptr

//go:nosplit
func sysAlloc(n uintptr, sysStat *uint64) unsafe.Pointer {
	p := sysReserve(nil, n)
	sysMap(p, n, sysStat)
	// println("sysAlloc, n=", n, " p=", p)
	return p
}

func sysUnused(v unsafe.Pointer, n uintptr) {
}

func sysUsed(v unsafe.Pointer, n uintptr) {
}

func sysHugePage(v unsafe.Pointer, n uintptr) {
}

//go:nosplit
func sysFree(v unsafe.Pointer, n uintptr, sysStat *uint64) {
	mSysStatDec(sysStat, n)
}

func sysFault(v unsafe.Pointer, n uintptr) {
}

func sysReserve(v unsafe.Pointer, n uintptr) unsafe.Pointer {
	if v != nil {
		return nil
	}

	p := memoryEnd
	memoryEnd += n
	return unsafe.Pointer(p)
}

func sysMap(v unsafe.Pointer, n uintptr, sysStat *uint64) {
	mSysStatInc(sysStat, n)
}
