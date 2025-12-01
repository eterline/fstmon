// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package hostdata

import (
	"bytes"
	"encoding/hex"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type HexID [16]byte

func (id HexID) String() string {
	return hex.EncodeToString(id[:])
}

/*
Hostnamectl – handle information from command: "hostnamectl"

	┌──────────────────────┬────────────────────────────────────────────────────────┐
	│ Field                │ Description                                            │
	├──────────────────────┼────────────────────────────────────────────────────────┤
	│ Hostname             │ The system hostname (e.g., "myserver.local")           │
	│ IconName             │ Device icon type (e.g., "computer-laptop",             │
	│                      │ "computer-server"), used in graphical UIs              │
	│ Chassis              │ Device chassis type: "laptop", "desktop", "server",    │
	│                      │ "vm", etc.                                             │
	│ MachineID            │ Unique machine identifier (hex), stored in             │
	│                      │ /etc/machine-id, persistent across reboots             │
	│ BootID               │ Unique identifier for the current system boot session  │
	│ Kernel               │ Linux kernel version (e.g., "5.15.0-105-generic")      │
	│ Architecture         │ CPU architecture: "x86_64", "aarch64", etc.            │
	│ HardwareVendor       │ Hardware manufacturer (e.g., "Dell Inc.")              │
	│ HardwareModel        │ Hardware model name (e.g., "ThinkPad T490")            │
	│ FirmwareVersion      │ BIOS/UEFI firmware version (e.g., "1.12.0")            │
	│ FirmwareDate         │ Release date of the firmware                           │
	│ FirmwareAge          │ Firmware age (calculated as the time since             │
	│                      │ the FirmwareDate)                                      │
	└──────────────────────┴────────────────────────────────────────────────────────┘
*/
type Hostnamectl struct {
	Hostname        string
	IconName        string
	Chassis         string
	MachineID       HexID
	BootID          HexID
	Kernel          string
	Architecture    string
	HardwareVendor  string
	HardwareModel   string
	FirmwareVersion string
	FirmwareDate    time.Time
	FirmwareAge     time.Duration
}

var (
	hnameCtl Hostnamectl
)

func init() {
	// pre starting reading information
	readHostnamectl(&hnameCtl)
}

// parses 'hostnamectl' output
func readHostnamectl(h *Hostnamectl) {
	data, err := exec.Command("hostnamectl").Output()
	if err != nil {
		return
	}

	lines := bytes.Split(data, []byte{'\n'})
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte{':'}, 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(string(parts[0]))
		val := parts[1]

		switch key {
		case "Static hostname":
			h.Hostname = bytesToString(val)
		case "Icon name":
			h.IconName = bytesToString(val)
		case "Chassis":
			h.Chassis = bytesToString(val)
		case "Machine ID":
			h.MachineID = parseHexVal(bytesToString(val))
		case "Boot ID":
			h.BootID = parseHexVal(bytesToString(val))
		case "Kernel":
			h.Kernel = bytesToString(val)
		case "Architecture":
			h.Architecture = bytesToString(val)
		case "Hardware Vendor":
			h.HardwareVendor = bytesToString(val)
		case "Hardware Model":
			h.HardwareModel = bytesToString(val)
		case "Firmware Version":
			h.FirmwareVersion = bytesToString(val)
		case "Firmware Date":
			value := bytesToString(val)
			h.FirmwareDate, _ = time.Parse("Mon 2006-01-02", value)
		case "Firmware Age":
			value := bytesToString(val)
			h.FirmwareAge = parseCustomDuration(value)
		}
	}
}

// GetHostnamectl return hostnamectl information
func GetHostnamectl() Hostnamectl {
	return hnameCtl
}

func removeNonASCII(s string) string {
	var result []rune
	for _, r := range s {
		if r <= 127 {
			result = append(result, r)
		}
	}
	return string(result)
}

func parseCustomDuration(s string) time.Duration {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`(\d+)\s*(y|year|years|month|months|d|day|days)`)
	matches := re.FindAllStringSubmatch(s, -1)

	var total time.Duration
	for _, match := range matches {
		numStr := match[1]
		unit := match[2]

		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}

		switch unit {
		case "y", "year", "years":
			total += time.Duration(num) * 365 * 24 * time.Hour
		case "month", "months":
			total += time.Duration(num) * 30 * 24 * time.Hour
		case "d", "day", "days":
			total += time.Duration(num) * 24 * time.Hour
		}
	}

	return total
}

func parseHexVal(s string) [16]byte {
	res := [16]byte{}
	id, err := hex.DecodeString(s)
	if err == nil {
		copy(res[:], id)
	}
	return res
}
