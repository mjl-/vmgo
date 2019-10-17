// +build !solo5hvt

package os

func (f *File) isFake() bool {
	return false
}
