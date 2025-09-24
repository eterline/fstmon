package secure

import (
	"log/slog"
	"strings"
)

type AllowedHostsFilter struct {
	ahf map[string]struct{}
}

func InitAllowedHostsFilter(host ...string) *AllowedHostsFilter {
	hosts := make(map[string]struct{}, len(host))
	for _, h := range host {
		hosts[strings.ToLower(h)] = struct{}{}
	}

	f := &AllowedHostsFilter{
		ahf: hosts,
	}

	if len(hosts) > 0 {
		slog.Warn("host filter enabled", "allow", f.allowedHosts())
	}
	return f
}

func (f *AllowedHostsFilter) allowedHosts() []string {
	a := make([]string, 0, len(f.ahf))
	for host := range f.ahf {
		a = append(a, host)
	}
	return a
}

func (f *AllowedHostsFilter) InAllowedHosts(host string) bool {
	if f.ahf != nil {
		return true
	}
	_, ok := f.ahf[strings.ToLower(host)]
	return ok
}
