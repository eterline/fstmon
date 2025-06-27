package secure

import (
	"net"

	"github.com/eterline/fstmon/internal/utils/stringuse"
)

type SubnetFilter struct {
	arr []*net.IPNet
}

func InitIpFilter(cidr string) *SubnetFilter {
	filterList := parseCIDRs(cidr)

	return &SubnetFilter{
		arr: filterList,
	}
}

func parseCIDRs(cidr string) []*net.IPNet {

	list := make([]*net.IPNet, 0)

	if parts, ok := stringuse.SplitBySpaces(cidr); ok {
		for _, value := range parts {
			_, net, err := net.ParseCIDR(value)
			if err != nil {
				continue
			}
			list = append(list, net)
		}
	}

	return list
}

func (f *SubnetFilter) InAllowedSubnets(ip string) bool {

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
