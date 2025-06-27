package secure

import (
	"github.com/eterline/fstmon/internal/utils/stringuse"
)

func parseAllowedHosts(hosts string) map[string]struct{} {

	allowed := make(map[string]struct{}, 0)

	if parts, ok := stringuse.SplitBySpaces(hosts); ok {
		for _, value := range parts {
			if value != "" {
				allowed[value] = struct{}{}
			}
		}
	}

	return allowed
}

type AllowedHostsFilter struct {
	allowHosts map[string]struct{}
}

func InitAllowedHostsFilter(host string) *AllowedHostsFilter {
	filterList := parseAllowedHosts(host)

	return &AllowedHostsFilter{
		allowHosts: filterList,
	}
}

func (f *AllowedHostsFilter) InAllowedHosts(host string) bool {

	if len(f.allowHosts) == 0 {
		return true
	}

	_, ok := f.allowHosts[host]
	return ok
}
