package usecase_test

import (
	"testing"

	"github.com/eterline/fstmon/internal/utils/usecase"
)

func Test_AvgVector(t *testing.T) {
	t.Run("int slice", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []int
			expected int
		}{
			{"empty slice", []int{}, 0},
			{"single element", []int{5}, 5},
			{"multiple elements", []int{1, 2, 3, 4}, 2}, // integer division
			{"negative numbers", []int{-1, -2, -3}, -2},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := usecase.AvgVector(tt.input)
				if got != tt.expected {
					t.Errorf("AvgVector(%v) = %v, want %v", tt.input, got, tt.expected)
				}
			})
		}
	})

	t.Run("float64 slice", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []float64
			expected float64
		}{
			{"empty slice", []float64{}, 0},
			{"single element", []float64{5.5}, 5.5},
			{"multiple elements", []float64{1.5, 2.5, 3.0}, 2.3333333333333335},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := usecase.AvgVector(tt.input)
				if got != tt.expected {
					t.Errorf("AvgVector(%v) = %v, want %v", tt.input, got, tt.expected)
				}
			})
		}
	})
}

func Test_AvgVectorFunc(t *testing.T) {
	type item struct {
		val int
	}

	t.Run("basic int", func(t *testing.T) {
		s := []item{{1}, {2}, {3}, {4}}
		got := usecase.AvgVectorFunc(s, func(idx int) int {
			return s[idx].val
		})
		expected := 2 // integer division
		if got != expected {
			t.Errorf("AvgVectorFunc = %v, want %v", got, expected)
		}
	})

	t.Run("float64", func(t *testing.T) {
		type fitem struct {
			v float64
		}
		s := []fitem{{1.5}, {2.5}, {3.0}}
		got := usecase.AvgVectorFunc(s, func(idx int) float64 {
			return s[idx].v
		})
		expected := (1.5 + 2.5 + 3.0) / 3
		if got != expected {
			t.Errorf("AvgVectorFunc = %v, want %v", got, expected)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		s := []item{}
		got := usecase.AvgVectorFunc(s, func(idx int) int {
			return 1
		})
		if got != 0 {
			t.Errorf("AvgVectorFunc(empty) = %v, want 0", got)
		}
	})

	t.Run("single element", func(t *testing.T) {
		s := []item{{42}}
		got := usecase.AvgVectorFunc(s, func(idx int) int {
			return s[idx].val
		})
		if got != 42 {
			t.Errorf("AvgVectorFunc(single) = %v, want 42", got)
		}
	})
}
