package syscall

// from fs_nacl.go

func Getcwd(buf []byte) (n int, err error) {
	// Force package os to default to the old algorithm using .. and directory reads.
	return 0, ENOSYS
}
