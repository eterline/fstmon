package ipfilter

import (
	"net"
	"regexp"
	"strings"
)

var splitRegex = regexp.MustCompile(`[,\s;]+`)

type IpFilter struct {
	arr []*net.IPNet
}

func InitIpFilter(cidr string) *IpFilter {
	filterList := parseCIDRs(cidr)

	return &IpFilter{
		arr: filterList,
	}
}

func parseCIDRs(cidr string) []*net.IPNet {
	parts := splitRegex.Split(strings.TrimSpace(cidr), -1)
	list := make([]*net.IPNet, 0)

	for _, value := range parts {
		_, net, err := net.ParseCIDR(value)
		if err != nil {
			continue
		}
		list = append(list, net)
	}

	return list
}

func (f *IpFilter) InAllowedCIDR(ip string) bool {

	if len(f.arr) == 0 {
		return true
	}

	ipEq := net.ParseIP(ip)
	if ipEq == nil {
		return false
	}

	for _, subnet := range f.arr {
		if subnet.Contains(ipEq) {
			return true
		}
	}

	return false
}
