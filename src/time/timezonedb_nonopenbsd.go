// +build !openbsd

package time

// Nop on non-openbsd.
func SetTimezoneDB(zipData []byte) {
}
