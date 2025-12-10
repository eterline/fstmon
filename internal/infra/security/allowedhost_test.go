// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package security_test

import (
	"strings"
	"testing"

	"github.com/eterline/fstmon/internal/infra/security"
)

func TestInitAllowedHostsFilter(t *testing.T) {

	filter := security.InitAllowedHostsFilter("example.com", "Test.com")
	if filter == nil {
		t.Fatal("Expected filter, got nil")
	}

	if len(filter.AllowedHosts()) != 2 {
		t.Fatalf("Expected 2 hosts, got %d", len(filter.AllowedHosts()))
	}

	if !filter.InAllowedHosts("EXAMPLE.COM") {
		t.Error("Expected EXAMPLE.COM to be allowed")
	}

	if !filter.InAllowedHosts("test.com") {
		t.Error("Expected test.com to be allowed")
	}
}

func TestInitAllowedHostsFilterEmpty(t *testing.T) {
	filter := security.InitAllowedHostsFilter()
	if filter != nil {
		t.Error("Expected nil filter for empty input")
	}
}

func TestInAllowedHosts(t *testing.T) {
	filter := security.InitAllowedHostsFilter("one.com", "two.com")

	// Разрешенные хосты
	if !filter.InAllowedHosts("one.com") {
		t.Error("Expected one.com to be allowed")
	}
	if !filter.InAllowedHosts("TWO.COM") {
		t.Error("Expected TWO.COM to be allowed (case-insensitive)")
	}

	// Не разрешенный хост
	if filter.InAllowedHosts("three.com") {
		t.Error("Expected three.com to NOT be allowed")
	}

	// Проверка поведения для nil фильтра
	var nilFilter *security.AllowedHostsFilter
	if !nilFilter.InAllowedHosts("any.com") {
		t.Error("Expected nil filter to allow any host")
	}
}

func TestAllowedHosts(t *testing.T) {
	hosts := []string{"a.com", "b.com", "c.com"}
	filter := security.InitAllowedHostsFilter(hosts...)
	got := filter.AllowedHosts()

	hostMap := make(map[string]bool)
	for _, h := range got {
		hostMap[h] = true
	}

	for _, h := range hosts {
		if !hostMap[strings.ToLower(h)] {
			t.Errorf("Expected host %s in AllowedHosts", h)
		}
	}
}
