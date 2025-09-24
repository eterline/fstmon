// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parseInterfaceCounters(t *testing.T) {
	data := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo:     392       6    0    0    0     0          0         0      392       6    0    0    0     0       0          0
enp5s0: 889517540  657028    0    0    0     0          0         0 27480511  189071    0    2    0     0       0          0
enp5s1: 889  657028    0    0    0     0          0         0 889  189071    0    2    0     0       0          0
`

	trg := InterfaceCounters{
		"lo": InterfaceCounter{
			Receive:  NetRX{392, 6, 0, 0, 0, 0, 0, 0},
			Transmit: NetTX{392, 6, 0, 0, 0, 0, 0, 0},
		},

		"enp5s0": InterfaceCounter{
			Receive:  NetRX{889517540, 657028, 0, 0, 0, 0, 0, 0},
			Transmit: NetTX{27480511, 189071, 0, 2, 0, 0, 0, 0},
		},

		"enp5s1": InterfaceCounter{
			Receive:  NetRX{889, 657028, 0, 0, 0, 0, 0, 0},
			Transmit: NetTX{889, 189071, 0, 2, 0, 0, 0, 0},
		},
	}

	fmt.Println(trg)

	ct, err := parseInterfaceCounters([]byte(data))
	if err != nil {
		t.Error(err)
	}

	fmt.Println(ct)

	if r := cmp.Diff(trg, ct); r != "" {
		t.Error(r)
	}
}
