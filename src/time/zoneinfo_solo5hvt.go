// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"syscall"
)

var zoneSources = []string{
	"/zoneinfo.zip", // Recognized in sys_solo5hvt.go.
}

func initLocal() {
	// consult $TZ to find the time zone to use.
	// nor $TZ or $TZ="" means use UTC.
	// $TZ="foo" means use /usr/share/zoneinfo/foo.

	tz, _ := syscall.Getenv("TZ")
	if tz != "" && tz != "UTC" {
		if z, err := loadLocation(tz, zoneSources); err == nil {
			localLoc = *z
			return
		}
	}

	// Fall back to UTC.
	localLoc.name = "UTC"
}
