package secure

import (
	"log/slog"
	"net/netip"

	"github.com/eterline/fstmon/pkg/netipuse"
)

type SubnetFilter struct {
	pool *netipuse.PoolIP
}

func NewSubnetFilter(cidr []string) *SubnetFilter {
	if len(cidr) < 1 {
		return &SubnetFilter{}
	}

	pBuild := netipuse.NewPoolIPBuilder()
	for _, sub := range cidr {
		err := pBuild.AddPrefixParseSubnet(sub)
		if err != nil {
			slog.Info("invalid filter subnet", "err", err.Error())
		}
	}

	pool, err := pBuild.PoolIP()
	if err != nil {
		slog.Error("failed to init subnet filter", "err", err.Error())
	}

	slog.Warn("subnet filter enabled", "allow", pool.Prefixes())
	return &SubnetFilter{
		pool: pool,
	}
}

func (f *SubnetFilter) InAllowedSubnets(ip netip.Addr) bool {
	if (f.pool == nil) || f.pool.Contains(ip) {
		return true
	}
	return ip.IsLoopback()
}
