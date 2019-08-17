// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

type rawConn struct{}

func newRawConn(fd *netFD) (*rawConn, error)        { return nil, notyet("newRawConn") }
func (c *rawConn) ok() bool                         { return false }
func (c *rawConn) Control(f func(uintptr)) error    { return notyet("rawConn.Control") }
func (c *rawConn) Read(f func(uintptr) bool) error  { return notyet("rawConn.Read") }
func (c *rawConn) Write(f func(uintptr) bool) error { return notyet("rawConn.Write") }

type rawListener struct {
	rawConn
}

func newRawListener(fd *netFD) (*rawListener, error)  { return nil, notyet("newRawListener") }
func (l *rawListener) Read(func(uintptr) bool) error  { return notyet("rawListener.Read") }
func (l *rawListener) Write(func(uintptr) bool) error { return notyet("rawListener.Write") }
