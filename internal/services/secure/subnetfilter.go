// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package secure

import (
	"fmt"
	"net/netip"

	"github.com/eterline/fstmon/pkg/netipuse"
)

type SubnetFilter struct {
	pool *netipuse.PoolIP
}

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

func (f *SubnetFilter) InAllowedSubnets(ip netip.Addr) bool {
	if (f.pool == nil) || f.pool.Contains(ip) {
		return true
	}
	return ip.IsLoopback()
}

func (f *SubnetFilter) AllowedList() []netip.Prefix {
	if f.pool == nil {
		return nil
	}
	return f.pool.Prefixes()
}
