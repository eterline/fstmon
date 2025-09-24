// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type ThermalValue struct {
	Current float64
	Minimum float64
	Maximum float64
}

func Temperatures(ctx context.Context) (map[string]ThermalValue, error) {

	data, err := exec.CommandContext(ctx, "sensors", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("error fetch temp sensors: %w", err)
	}

	var parsed map[string]interface{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, fmt.Errorf("error parsing temperatures: %w", err)
	}

	return extractThermalValues(parsed), nil
}

func extractThermalValues(data map[string]interface{}) map[string]ThermalValue {
	result := make(map[string]ThermalValue)

	var walk func(prefix string, v interface{})
	walk = func(prefix string, v interface{}) {
		switch val := v.(type) {
		case map[string]interface{}:
			var tv ThermalValue
			var found bool

			for k, v2 := range val {
				lkey := strings.ToLower(k)

				if strings.Contains(lkey, "input") {
					if f, ok := v2.(float64); ok {
						tv.Current = f
						found = true
					}
				}

				if strings.Contains(lkey, "min") {
					if f, ok := v2.(float64); ok {
						tv.Minimum = f
					}
				}

				if strings.Contains(lkey, "max") {
					if f, ok := v2.(float64); ok {
						tv.Maximum = f
					}
				}
			}

			if found {
				result[prefix] = tv
			}

			for k, v2 := range val {
				k = strings.ToLower(k)
				k = strings.ReplaceAll(k, " ", "-")
				walk(prefix+"_"+k, v2)
			}
		}
	}

	for k, v := range data {
		walk(k, v)
	}

	return result
}
