package os

// Full path to contents, eg "/etc/resolv.conf".
var fakeFiles = map[string][]byte{}

// Register a path with given contents.
// A later os.Open() on the given path returns a file from which the original contents can be read.
func AddFile(path string, data []byte) {
	fakeFiles[path] = data
}
