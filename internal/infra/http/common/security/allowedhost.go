// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package security

import (
	"strings"
)

type AllowedHostsFilter struct {
	ahf map[string]struct{}
}

func InitAllowedHostsFilter(host ...string) *AllowedHostsFilter {
	if len(host) == 0 {
		return nil
	}

	hosts := make(map[string]struct{}, len(host))
	for _, h := range host {
		hosts[strings.ToLower(h)] = struct{}{}
	}

	f := &AllowedHostsFilter{
		ahf: hosts,
	}

	return f
}

func (f *AllowedHostsFilter) AllowedHosts() []string {
	a := make([]string, 0, len(f.ahf))
	for host := range f.ahf {
		a = append(a, host)
	}
	return a
}

func (f *AllowedHostsFilter) InAllowedHosts(host string) bool {
	if f == nil {
		return true
	}
	_, ok := f.ahf[strings.ToLower(host)]
	return ok
}
