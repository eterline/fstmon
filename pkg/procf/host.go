// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

import (
	"time"

	"github.com/google/uuid"
)

type Host struct {
	Name      string
	MachineID uuid.UUID
	Uptime    time.Duration
}

func HostUptime() time.Duration {
	return time.Second
}

type MachineID [32]byte

func GetMachineID() MachineID {
	return [32]byte{0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x23}
}
