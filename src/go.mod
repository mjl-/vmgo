module std

go 1.12

require (
	github.com/google/netstack v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7
	golang.org/x/sys v0.0.0-20190529130038-5219a1e1c5f8 // indirect
	golang.org/x/text v0.3.2 // indirect
)

replace github.com/google/netstack => github.com/mjl-/netstack v0.0.0-20190823123829-92d57aa60fe3

// replace github.com/google/netstack => ../../netstack
