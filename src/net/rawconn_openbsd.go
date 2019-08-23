// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "syscall"

type rawConn struct{}

func newRawConn(fd *netFD) (*rawConn, error)        { return nil, syscall.ENOSYS }
func (c *rawConn) ok() bool                         { return false }
func (c *rawConn) Control(f func(uintptr)) error    { return syscall.ENOSYS }
func (c *rawConn) Read(f func(uintptr) bool) error  { return syscall.ENOSYS }
func (c *rawConn) Write(f func(uintptr) bool) error { return syscall.ENOSYS }

type rawListener struct {
	rawConn
}

func newRawListener(fd *netFD) (*rawListener, error)  { return nil, syscall.ENOSYS }
func (l *rawListener) Read(func(uintptr) bool) error  { return syscall.ENOSYS }
func (l *rawListener) Write(func(uintptr) bool) error { return syscall.ENOSYS }
