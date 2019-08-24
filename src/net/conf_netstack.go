// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build netstack

package net

import (
	"errors"
	"sync"
	"time"
)

type conf struct {
	resolv *dnsConfig
}

var (
	confOnce sync.Once // guards init of confVal via initConfVal
	confVal  = &conf{}
)

// systemConf returns the machine's network configuration.
func systemConf() *conf {
	confOnce.Do(initConfVal)
	return confVal
}

func initConfVal() {
	confVal.resolv = &dnsConfig{
		ndots:         1,
		timeout:       5 * time.Second,
		attempts:      2,
		singleRequest: true,
	}

	if len(netstackNameservers) == 0 {
		confVal.resolv.err = errors.New("no nameservers configured")
		return
	}
	confVal.resolv.servers = netstackNameservers
}

// canUseCgo reports whether calling cgo functions is allowed
// for non-hostname lookups.
func (c *conf) canUseCgo() bool {
	return false
}

// hostLookupOrder determines which strategy to use to resolve hostname.
// The provided Resolver is optional. nil means to not consider its options.
func (c *conf) hostLookupOrder(r *Resolver, hostname string) (ret hostLookupOrder) {
	return hostLookupDNS
}
