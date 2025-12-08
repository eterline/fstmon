// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package security_test

import (
	"net/http"
	"testing"

	"github.com/eterline/fstmon/internal/infra/security"
)

func TestParseXRealIP(t *testing.T) {
	tests := []struct {
		name   string
		header string
		wantOk bool
		wantIP string
	}{
		{"empty", "", false, ""},
		{"valid", "192.168.1.1", true, "192.168.1.1"},
		{"invalid", "invalid-ip", false, ""},
	}

	for _, tt := range tests {
		h := http.Header{}
		h.Set("X-Real-IP", tt.header)
		ip, ok := security.ParseXRealIP(h)
		if ok != tt.wantOk {
			t.Errorf("%s: got ok %v, want %v", tt.name, ok, tt.wantOk)
		}
		if ok && ip.String() != tt.wantIP {
			t.Errorf("%s: got ip %v, want %v", tt.name, ip.String(), tt.wantIP)
		}
	}
}

func TestParseXForwardedFor(t *testing.T) {
	tests := []struct {
		name   string
		header string
		wantOk bool
		wantIP string
	}{
		{"empty", "", false, ""},
		{"single", "10.0.0.1", true, "10.0.0.1"},
		{"multiple", "10.0.0.1, 10.0.0.2", true, "10.0.0.1"},
		{"invalid", "invalid-ip", false, ""},
	}

	for _, tt := range tests {
		h := http.Header{}
		h.Set("X-Forwarded-For", tt.header)
		ip, ok := security.ParseXForwardedFor(h)
		if ok != tt.wantOk {
			t.Errorf("%s: got ok %v, want %v", tt.name, ok, tt.wantOk)
		}
		if ok && ip.String() != tt.wantIP {
			t.Errorf("%s: got ip %v, want %v", tt.name, ip.String(), tt.wantIP)
		}
	}
}

func TestParseForwarded(t *testing.T) {
	tests := []struct {
		name   string
		header string
		wantOk bool
		wantIP string
	}{
		{"empty", "", false, ""},
		{"valid", "for=192.0.2.60", true, "192.0.2.60"},
		{"quoted", "for=\"[2001:db8::1]\"", true, "2001:db8::1"},
		{"multiple", "for=192.0.2.60, for=198.51.100.17", true, "192.0.2.60"},
		{"invalid", "for=invalid-ip", false, ""},
	}

	for _, tt := range tests {
		h := http.Header{}
		h.Set("Forwarded", tt.header)
		ip, ok := security.ParseForwarded(h)
		if ok != tt.wantOk {
			t.Errorf("%s: got ok %v, want %v", tt.name, ok, tt.wantOk)
		}
		if ok && ip.String() != tt.wantIP {
			t.Errorf("%s: got ip %v, want %v", tt.name, ip.String(), tt.wantIP)
		}
	}
}

func TestRemote(t *testing.T) {
	validIP := "192.168.1.100:12345"
	invalidIP := "not-an-ip"

	r := &http.Request{RemoteAddr: validIP}
	ip, err := security.Remote(r)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if ip.Addr().String() != "192.168.1.100" {
		t.Errorf("expected 192.168.1.100, got %v", ip)
	}

	r.RemoteAddr = invalidIP
	_, err = security.Remote(r)
	if err == nil {
		t.Errorf("expected error for invalid remote addr")
	}
}

func TestExtractIP_Headers(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		wantIP  string
		wantOk  bool
	}{
		{"xrealip", map[string]string{"X-Real-IP": "1.1.1.1"}, "1.1.1.1", true},
		{"xforwardedfor", map[string]string{"X-Forwarded-For": "2.2.2.2"}, "2.2.2.2", true},
		{"forwarded", map[string]string{"Forwarded": "for=3.3.3.3"}, "3.3.3.3", true},
		{"remote", map[string]string{}, "4.4.4.4", true},
	}

	// TODO: fix errors parsing
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{Header: http.Header{}}
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			if tt.name == "remote" {
				req.RemoteAddr = "4.4.4.4:1234"
			} else {
				req.RemoteAddr = "invalid:addr"
			}

			extractor := security.NewIpExtractor(true)
			ip, _, err := extractor.ExtractIP(req)
			if tt.wantOk && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantOk && ip.String() != tt.wantIP {
				t.Errorf("expected %v, got %v", tt.wantIP, ip)
			}
		})
	}
}

func TestExtractIP_NoHeaders(t *testing.T) {
	req := &http.Request{RemoteAddr: "5.5.5.5:1234"}
	extractor := security.NewIpExtractor(false)
	ip, _, err := extractor.ExtractIP(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ip.String() != "5.5.5.5" {
		t.Errorf("expected 5.5.5.5, got %v", ip)
	}
}
