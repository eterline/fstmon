package security

import (
	"fmt"
	"net/netip"

	"github.com/eterline/fstmon/pkg/netipuse"
)

/*
SubnetFilter – represents a filter for allowed subnets.

	Holds a pool of allowed IP prefixes. Provides methods to check
	if a given IP is allowed and to retrieve the list of allowed prefixes.
*/
type SubnetFilter struct {
	pool *netipuse.PoolIP
}

/*
NewSubnetFilter – creates a new SubnetFilter from a list of CIDR strings.

	Attempts to parse each CIDR and builds a pool of allowed subnets.
	Returns the filter and an aggregated error if any of the subnets failed to parse.
*/
func NewSubnetFilter(cidr []string) (*SubnetFilter, error) {
	if len(cidr) < 1 {
		return &SubnetFilter{}, nil
	}

	errs := []error{}

	pBuild := netipuse.NewPoolIPBuilder()
	for _, sub := range cidr {
		err := pBuild.AddPrefixParseSubnet(sub)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid filter subnet: %w", err))
		}
	}

	pool, err := pBuild.PoolIP()
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to init subnet filter: %w", err))
	}

	filter := &SubnetFilter{
		pool: pool,
	}

	if len(errs) == 0 {
		return filter, nil
	}

	return filter, fmt.Errorf("subnet filter errors: %v", errs)
}

/*
InAllowedSubnets – checks whether the given IP address is within the allowed subnets.

	Returns true if:
	- The pool is nil (no restriction),
	- The IP is contained in the pool,
	- The IP is a loopback address.
*/
func (f *SubnetFilter) InAllowedSubnets(ip netip.Addr) bool {
	if (f.pool == nil) || f.pool.Contains(ip) {
		return true
	}
	return ip.IsLoopback()
}

/*
AllowedList – returns the list of allowed prefixes as netip.Prefix slices.

	Returns nil if no pool is defined.
*/
func (f *SubnetFilter) AllowedList() ([]netip.Prefix, bool) {
	if f.pool == nil {
		return nil, false
	}
	pfxs := f.pool.Prefixes()
	return pfxs, len(pfxs) > 0
}
