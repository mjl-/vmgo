package os

import (
	"syscall"
)

func removeAll(path string) error {
	return syscall.ENOTSUP
}
