// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package security_test

import (
	"net/netip"
	"testing"

	"github.com/eterline/fstmon/internal/infra/http/common/security"
)

func TestNewSubnetFilter_ValidAndInvalid(t *testing.T) {
	cidrs := []string{"192.168.1.0/24", "10.0.0.0/8", "invalid-cidr"}
	filter, err := security.NewSubnetFilter(cidrs)
	if filter == nil {
		t.Fatal("Expected filter, got nil")
	}

	if err == nil {
		t.Error("Expected error due to invalid CIDR, got nil")
	}

	if !filter.InAllowedSubnets(netip.MustParseAddr("192.168.1.1")) {
		t.Error("192.168.1.1 should be allowed")
	}
	if !filter.InAllowedSubnets(netip.MustParseAddr("10.0.0.1")) {
		t.Error("10.0.0.1 should be allowed")
	}
}

func TestNewSubnetFilter_Empty(t *testing.T) {
	filter, err := security.NewSubnetFilter([]string{})
	if err != nil {
		t.Errorf("Expected no error for empty CIDR, got %v", err)
	}
	if filter == nil {
		t.Fatal("Expected filter, got nil")
	}
}

func TestInAllowedSubnets(t *testing.T) {
	cidrs := []string{"192.168.0.0/16"}
	filter, err := security.NewSubnetFilter(cidrs)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	tests := []struct {
		ip      string
		allowed bool
	}{
		{"192.168.1.1", true},
		{"10.0.0.1", false},
		{"127.0.0.1", true},
	}

	for _, tt := range tests {
		addr := netip.MustParseAddr(tt.ip)
		if got := filter.InAllowedSubnets(addr); got != tt.allowed {
			t.Errorf("IP %v: expected %v, got %v", tt.ip, tt.allowed, got)
		}
	}
}

func TestInAllowedSubnets_NilPool(t *testing.T) {
	filter := &security.SubnetFilter{}
	tests := []string{"192.168.1.1", "10.0.0.1", "127.0.0.1"}

	for _, ipStr := range tests {
		ip := netip.MustParseAddr(ipStr)
		if !filter.InAllowedSubnets(ip) {
			t.Errorf("Expected %v to be allowed with nil pool", ipStr)
		}
	}
}

func TestAllowedList(t *testing.T) {
	cidrs := []string{"192.168.0.0/16", "10.0.0.0/8"}
	filter, err := security.NewSubnetFilter(cidrs)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	prefixes := filter.AllowedList()
	if prefixes == nil || len(prefixes) == 0 {
		t.Fatal("Expected non-empty prefix list")
	}

	found := make(map[string]bool)
	for _, p := range prefixes {
		found[p.String()] = true
	}

	for _, cidr := range cidrs {
		if !found[cidr] {
			t.Errorf("Expected CIDR %v in AllowedList", cidr)
		}
	}
}

func TestAllowedList_NilPool(t *testing.T) {
	filter := &security.SubnetFilter{}
	if list := filter.AllowedList(); list != nil {
		t.Errorf("Expected nil AllowedList for nil pool, got %v", list)
	}
}
