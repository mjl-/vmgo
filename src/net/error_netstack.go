// +build netstack

package net

import (
	"errors"

	"github.com/google/netstack/tcpip"
)

// todo: wrap the error while keeping the original, improves isConnError

func wrapNetstackError(err *tcpip.Error) error {
	if err == nil {
		return nil
	}
	return errors.New(err.String())
}

func isConnError(err error) bool {
	switch err.Error() {
	case tcpip.ErrConnectionReset.String(), tcpip.ErrConnectionAborted.String():
		return true
	}
	return false
}
