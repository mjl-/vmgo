// +build netstack

package net

func maxListenerBacklog() int {
	// currently hardcoded in gonet
	return 10
}
