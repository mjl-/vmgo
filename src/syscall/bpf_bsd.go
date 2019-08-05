// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build dragonfly freebsd netbsd openbsd

// Berkeley packet filter for BSD variants

package syscall

// Deprecated: Use golang.org/x/net/bpf instead.
func BpfStmt(code, k int) *BpfInsn {
	return &BpfInsn{Code: uint16(code), K: uint32(k)}
}

// Deprecated: Use golang.org/x/net/bpf instead.
func BpfJump(code, k, jt, jf int) *BpfInsn {
	return &BpfInsn{Code: uint16(code), Jt: uint8(jt), Jf: uint8(jf), K: uint32(k)}
}

// Deprecated: Use golang.org/x/net/bpf instead.
func BpfBuflen(fd int) (int, error) {
	return 0, ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func SetBpfBuflen(fd, l int) (int, error) {
	return 0, ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func BpfDatalink(fd int) (int, error) {
	return 0, ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func SetBpfDatalink(fd, t int) (int, error) {
	return 0, ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func SetBpfPromisc(fd, m int) error {
	return ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func FlushBpf(fd int) error {
	return ENOSYS
}

type ivalue struct {
	name  [IFNAMSIZ]byte
	value int16
}

// Deprecated: Use golang.org/x/net/bpf instead.
func BpfInterface(fd int, name string) (string, error) {
	return "", ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func SetBpfInterface(fd int, name string) error {
	return ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func BpfTimeout(fd int) (*Timeval, error) {
	return nil, ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func SetBpfTimeout(fd int, tv *Timeval) error {
	return ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func BpfStats(fd int) (*BpfStat, error) {
	return nil, ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func SetBpfImmediate(fd, m int) error {
	return ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func SetBpf(fd int, i []BpfInsn) error {
	return ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func CheckBpfVersion(fd int) error {
	return ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func BpfHeadercmpl(fd int) (int, error) {
	return 0, ENOSYS
}

// Deprecated: Use golang.org/x/net/bpf instead.
func SetBpfHeadercmpl(fd, f int) error {
	return ENOSYS
}
