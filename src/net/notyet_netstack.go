// +build netstack

package net

import "errors"

func notyet(s string) error {
	return errors.New("net: " + s + " not yet")
}
