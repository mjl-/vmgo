package net

import (
	"errors"

	"github.com/google/netstack/tcpip"
)

func wrapNetstackError(err *tcpip.Error) error {
	if err == nil {
		return nil
	}
	return errors.New(err.String())
}
