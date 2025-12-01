// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package fastparse_test

import (
	"testing"

	"github.com/eterline/fstmon/internal/utils/fastparse"
)

func BenchmarkParseFloat(b *testing.B) {

	str := []byte("11.32")

	for n := 0; n < b.N; n++ {
		_, _ = fastparse.ParseFloat[float64](str)
	}
}

func almostEqual[F float32 | float64](a, b F, epsilon F) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= epsilon
}

func TestParseFloatStrictFormat(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want32    float32
		want64    float64
		expectErr bool
	}{
		// Valid dataset
		{"integer positive", "123", 123, 123, false},
		{"integer negative", "-456", -456, -456, false},
		{"integer positive with +", "+789", 789, 789, false},
		{"float simple", "3.14", 3.14, 3.14, false},
		{"float negative", "-2.718", -2.718, -2.718, false},
		{"float zero integer part", "0.123", 0.123, 0.123, false},
		{"float zero fractional part", "123.0", 123, 123, false},

		// Invalid dataset
		{"empty string", "", 0, 0, true},
		{"double dot", "1.2.3", 0, 0, true},
		{"invalid chars", "12a34", 0, 0, true},
		{"only dot", ".", 0, 0, true},
		{"sign only", "-", 0, 0, true},
		{"trailing letters", "123abc", 0, 0, true},
		{"leading letters", "abc123", 0, 0, true},
		{"space inside", "12 3", 0, 0, true},
		{"space at start", " 123", 0, 0, true},
		{"space at end", "123 ", 0, 0, true},
	}

	const epsilon32 float32 = 1e-6
	const epsilon64 float64 = 1e-12

	for _, tt := range tests {
		t.Run("float32 "+tt.name, func(t *testing.T) {
			got, err := fastparse.StringParseFloat[float32](tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("StringParseFloat[float32](%q) error = %v, wantErr %v", tt.input, err, tt.expectErr)
			}
			if !tt.expectErr && !almostEqual(got, float32(tt.want32), epsilon32) {
				t.Errorf("StringParseFloat[float32](%q) = %v, want %v", tt.input, got, tt.want32)
			}
		})

		t.Run("float64 "+tt.name, func(t *testing.T) {
			got, err := fastparse.StringParseFloat[float64](tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("StringParseFloat[float64](%q) error = %v, wantErr %v", tt.input, err, tt.expectErr)
			}
			if !tt.expectErr && !almostEqual(got, tt.want64, epsilon64) {
				t.Errorf("StringParseFloat[float64](%q) = %v, want %v", tt.input, got, tt.want64)
			}
		})
	}
}

func TestParseFloatDirectStrict(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		want32    float32
		want64    float64
		expectErr bool
	}{
		{"integer", []byte("42"), 42, 42, false},
		{"float", []byte("0.99"), 0.99, 0.99, false},
		{"negative float", []byte("-1.5"), -1.5, -1.5, false},
		{"invalid char", []byte("1a"), 0, 0, true},
		{"empty", []byte(""), 0, 0, true},
	}

	const epsilon32 float32 = 1e-6
	const epsilon64 float64 = 1e-12

	for _, tt := range tests {
		t.Run("float32 "+tt.name, func(t *testing.T) {
			got, err := fastparse.ParseFloat[float32](tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("ParseFloat[float32](%q) error = %v, wantErr %v", tt.input, err, tt.expectErr)
			}
			if !tt.expectErr && !almostEqual(got, float32(tt.want32), epsilon32) {
				t.Errorf("ParseFloat[float32](%q) = %v, want %v", tt.input, got, tt.want32)
			}
		})
		t.Run("float64 "+tt.name, func(t *testing.T) {
			got, err := fastparse.ParseFloat[float64](tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("ParseFloat[float64](%q) error = %v, wantErr %v", tt.input, err, tt.expectErr)
			}
			if !tt.expectErr && !almostEqual(got, tt.want64, epsilon64) {
				t.Errorf("ParseFloat[float64](%q) = %v, want %v", tt.input, got, tt.want64)
			}
		})
	}
}
