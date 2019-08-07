// +build !openbsd

package os

var printOpen = false

func PrintOpen(v bool) {
}

func (f *File) isFake() bool {
	return false
}
