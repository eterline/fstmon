// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

/*
Process holds detailed information about proc

	┌─────────────┬────────────────────────────────────────────────────────┐
	│ Value       │ Description                                            │
	├─────────────┼────────────────────────────────────────────────────────┤
	│ PID         │ Process ID                                             │
	│ PPID        │ Parent process ID                                      │
	│ Cpu         │ Cpu usage                                              │
	│ Memory      │ Memory usage                                           │
	│ User        │ Executed by username                                   │
	│ Command     │ Command name                                           │
	│ Class       │                                                        │
	│ State       │ Indicating process state                               │
	│ FullCommand │ Full exec command                                      │
	└─────────────┴────────────────────────────────────────────────────────┘
*/
type Process struct {
	PID         int32   `json:"pid"`
	PPID        int32   `json:"ppid"`
	Cpu         float64 `json:"cpu"`
	Memory      float64 `json:"mem"`
	User        string  `json:"user"`
	Command     string  `json:"command"`
	Class       string  `json:"class"`
	State       string  `json:"state"`
	FullCommand string  `json:"full_command"`
}

func SysProcessList() ([]Process, error) {

	output, err := exec.Command("ps", "-eo", psCmdFieldsArg).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'ps' command: %v", err)
	}

	var (
		parts     = strings.Split(string(output), "\n")
		procCount = len(parts) - 1
		procArr   = make([]Process, procCount-1)
	)

	for i, part := range parts[1:] {
		if part == "" {
			break
		}

		data := []byte(part)

		proc := Process{
			PID:         bytesToInt32(bytes.TrimSpace(data[:10])),
			PPID:        bytesToInt32(bytes.TrimSpace(data[10:21])),
			Cpu:         bytesToFloat64(bytes.TrimSpace(data[21:26])),
			Memory:      bytesToFloat64(bytes.TrimSpace(data[26:33])),
			User:        bytesToString(bytes.TrimSpace(data[33:50])),
			Command:     bytesToString(data[50:71]),
			Class:       bytesToString(bytes.TrimSpace(data[71:78])),
			State:       bytesToString(bytes.TrimSpace(data[78:84])),
			FullCommand: bytesToString(data[84:]),
		}

		procArr[i] = proc
	}

	return procArr, nil
}

/*
SelfStatus holds detailed information about this program process

	┌─────────────────────────────┬───────────────────────────────────────────────────────────────┐
	│ Field                       │ Description                                                   │
	├─────────────────────────────┼───────────────────────────────────────────────────────────────┤
	│ CapAmbMask                  │ Ambient capabilities mask (bitmask of permitted capabilities) │
	│ CapBndMask                  │ Capability bounding set mask (limits capabilities for thread) │
	│ CapEffMask                  │ Effective capabilities mask (currently effective capabilities)│
	│ CapInhMask                  │ Inheritable capabilities mask (inheritable by child processes)│
	│ CapPrmMask                  │ Permitted capabilities mask                                   │
	│ CoreDumping                 │ Core dumping flag (1 if core dump allowed, 0 otherwise)       │
	│ CpusAllowedMaskVec          │ CPUs allowed mask (bitmask vector representing allowed CPUs)  │
	│ CpusAllowedList             │ CPUs allowed list (string listing allowed CPUs)               │
	│ FDSize                      │ File descriptor limit (max open files for process)            │
	│ Gid                         │ Real, effective, saved set, and filesystem group IDs (list)   │
	│ Groups                      │ Supplementary group IDs (list)                                │
	│ HugetlbPages                │ Number of huge pages used                                     │
	│ Kthread                     │ Kernel thread flag (1 if kernel thread, else 0)               │
	│ MemsAllowed                 │ Allowed NUMA memory nodes (list of strings)                   │
	│ MemsAllowedList             │ Allowed memory nodes (string list)                            │
	│ NSpgid                      │ PID in PID namespace (process group ID)                       │
	│ NSpid                       │ PID in PID namespace (process ID)                             │
	│ NSsid                       │ Session ID in namespace                                       │
	│ NStgid                      │ Thread group ID in namespace                                  │
	│ Name                        │ Command name (process name)                                   │
	│ Ngid                        │ Thread group ID (kernel thread group ID)                      │
	│ NoNewPrivs                  │ No new privileges flag (1 if set, else 0)                     │
	│ PPid                        │ Parent process ID                                             │
	│ Pid                         │ Process ID                                                    │
	│ RssAnon                     │ Anonymous RSS (resident set size)                             │
	│ RssFile                     │ File-backed RSS                                               │
	│ RssShmem                    │ Shared memory RSS                                             │
	│ Seccomp                     │ Seccomp mode (0=no seccomp, 1=strict, 2=filter)               │
	│ SeccompFilters              │ Number of seccomp filters applied                             │
	│ ShdPnd                      │ Number of pages pending writeback                             │
	│ SigBlk                      │ Signal mask of blocked signals (bitmask)                      │
	│ SigCgt                      │ Signal mask of caught signals (bitmask)                       │
	│ SigIgn                      │ Signal mask of ignored signals (bitmask)                      │
	│ SigPnd                      │ Signal mask of pending signals (bitmask)                      │
	│ SigQ                        │ Queue of signals pending delivery (string)                    │
	│ SpeculationIndirectBranch   │ Status of speculation indirect branch mitigation (string)     │
	│ SpeculationStoreBypass      │ Status of speculation store bypass mitigation (string)        │
	│ State                       │ Process state (e.g. R=running, S=sleeping)                    │
	│ THPEnabled                  │ Transparent Huge Pages enabled (1=yes, 0=no)                  │
	│ Tgid                        │ Thread group ID (main thread’s PID)                           │
	│ Threads                     │ Number of threads in thread group                             │
	│ TracerPid                   │ PID of process tracing this process (0 if none)               │
	│ Uid                         │ Real, effective, saved set, and filesystem user IDs (list)    │
	│ VmData                      │ Size of data segment (in kB)                                  │
	│ VmExe                       │ Size of text segment (in kB)                                  │
	│ VmHWM                       │ Peak resident set size ("high water mark") in kB              │
	│ VmLck                       │ Locked memory size (in kB)                                    │
	│ VmLib                       │ Shared library code size (in kB)                              │
	│ VmPTE                       │ Page table entries size (in kB)                               │
	│ VmPeak                      │ Peak virtual memory size (in kB)                              │
	│ VmPin                       │ Pinned memory size (in kB)                                    │
	│ VmRSS                       │ Resident Set Size (in kB)                                     │
	│ VmSize                      │ Virtual memory size (in kB)                                   │
	│ VmStk                       │ Stack size (in kB)                                            │
	│ VmSwap                      │ Swapped-out virtual memory size (in kB)                       │
	│ NonvoluntaryCtxtSwitches    │ Number of forced context switches                             │
	│ UntagMask                   │ Untag mask (architecture-specific, related to hardware tags)  │
	│ VoluntaryCtxtSwitches       │ Number of voluntary context switches                          │
	│ X86ThreadFeatures           │ List of CPU thread features enabled (strings)                 │
	│ X86ThreadFeaturesLocked     │ List of locked CPU thread features (strings)                  │
	└─────────────────────────────┴───────────────────────────────────────────────────────────────┘
*/
type SelfStatus struct {
	CapAmbMask                uint64
	CapBndMask                uint64
	CapEffMask                uint64
	CapInhMask                uint64
	CapPrmMask                uint64
	CoreDumping               uint32
	CpusAllowedMaskVec        []uint8
	CpusAllowedList           string
	FDSize                    uint32
	Gid                       []uint16
	Groups                    []uint16
	HugetlbPages              uint64
	Kthread                   uint64
	MemsAllowed               []string
	MemsAllowedList           string
	NSpgid                    int64
	NSpid                     int64
	NSsid                     int64
	NStgid                    int64
	Name                      string
	Ngid                      int64
	NoNewPrivs                int64
	PPid                      int32
	Pid                       int32
	RssAnon                   uint64
	RssFile                   uint64
	RssShmem                  uint64
	Seccomp                   uint64
	SeccompFilters            uint64
	ShdPnd                    uint64
	SigBlk                    uint64
	SigCgt                    uint64
	SigIgn                    uint64
	SigPnd                    uint64
	SigQ                      string
	SpeculationIndirectBranch string
	SpeculationStoreBypass    string
	State                     string
	THPEnabled                uint64
	Tgid                      uint64
	Threads                   uint64
	TracerPid                 int32
	Uid                       []uint16
	VmData                    uint64
	VmExe                     uint64
	VmHWM                     uint64
	VmLck                     uint64
	VmLib                     uint64
	VmPTE                     uint64
	VmPeak                    uint64
	VmPin                     uint64
	VmRSS                     uint64
	VmSize                    uint64
	VmStk                     uint64
	VmSwap                    uint64
	NonvoluntaryCtxtSwitches  uint64
	UntagMask                 uint64
	VoluntaryCtxtSwitches     uint64
	X86ThreadFeatures         []string
	X86ThreadFeaturesLocked   []string
}

func FetchSelfStatus() (SelfStatus, error) {
	data, err := procSelfStatus.Data()
	if err != nil {
		return SelfStatus{}, err
	}

	prs := NewFileDataSetParser(
		data, []byte{'\n'},
		[]byte{':', '\t'},
		0, 1,
	)

	self := SelfStatus{
		CapAmbMask: byteMask(prs.Param("CapAmb")),
		CapBndMask: byteMask(prs.Param("CapBnd")),
		CapEffMask: byteMask(prs.Param("CapEff")),
		CapInhMask: byteMask(prs.Param("CapInh")),
		CapPrmMask: byteMask(prs.Param("CapPrm")),

		CoreDumping:        prs.Param("CoreDumping").Uint32(),
		CpusAllowedMaskVec: parseHexToBytes(prs.Param("Cpus_allowed")),
		CpusAllowedList:    prs.Param("Cpus_allowed_list").String(),
		FDSize:             prs.Param("FDSize").Uint32(),
		Gid:                gidList(prs.Param("Gid"), '\t'),
		Groups:             gidList(prs.Param("Groups"), ' '),
		HugetlbPages:       parseKbSize(prs.Param("HugetlbPages")),
		Kthread:            prs.Param("Kthread").Uint64(),
		MemsAllowed:        prs.Param("Mems_allowed").StringList(","),
		MemsAllowedList:    prs.Param("Mems_allowed_list").String(),
		NSpgid:             prs.Param("NSpgid").Int64(),
		NSpid:              prs.Param("NSpid").Int64(),
		NStgid:             prs.Param("NStgid").Int64(),
		Name:               prs.Param("Name").String(),
		Ngid:               prs.Param("Ngid").Int64(),
		NoNewPrivs:         prs.Param("NoNewPrivs").Int64(),
		PPid:               prs.Param("PPid").Int32(),
		Pid:                prs.Param("Pid").Int32(),
		RssAnon:            parseKbSize(prs.Param("RssAnon")),
		RssFile:            parseKbSize(prs.Param("RssFile")),
		RssShmem:           parseKbSize(prs.Param("RssShmem")),
		Seccomp:            prs.Param("Seccomp").Uint64(),
		SeccompFilters:     prs.Param("Seccomp_filters").Uint64(),
		ShdPnd:             byteMask(prs.Param("ShdPnd")),

		SigBlk: byteMask(prs.Param("SigBlk")),
		SigCgt: byteMask(prs.Param("SigCgt")),
		SigIgn: byteMask(prs.Param("SigIgn")),
		SigPnd: byteMask(prs.Param("SigPnd")),
		SigQ:   prs.Param("SigQ").String(),

		SpeculationIndirectBranch: prs.Param("SpeculationIndirectBranch").String(),
		SpeculationStoreBypass:    prs.Param("Speculation_Store_Bypass").String(),
		State:                     prs.Param("State").String(),
		THPEnabled:                prs.Param("THP_enabled").Uint64(),
		Tgid:                      prs.Param("Tgid").Uint64(),
		Threads:                   prs.Param("Threads").Uint64(),
		TracerPid:                 prs.Param("TracerPid").Int32(),
		Uid:                       gidList(prs.Param("Uid"), '\t'),

		VmData: parseKbSize(prs.Param("VmData")),
		VmExe:  parseKbSize(prs.Param("VmExe")),
		VmHWM:  parseKbSize(prs.Param("VmHWM")),
		VmLck:  parseKbSize(prs.Param("VmLck")),
		VmLib:  parseKbSize(prs.Param("VmLib")),
		VmPTE:  parseKbSize(prs.Param("VmPTE")),
		VmPeak: parseKbSize(prs.Param("VmPeak")),
		VmPin:  parseKbSize(prs.Param("VmPin")),
		VmRSS:  parseKbSize(prs.Param("VmRSS")),
		VmSize: parseKbSize(prs.Param("VmSize")),
		VmStk:  parseKbSize(prs.Param("VmStk")),
		VmSwap: parseKbSize(prs.Param("VmSwap")),

		NonvoluntaryCtxtSwitches: prs.Param("nonvoluntary_ctxt_switches").Uint64(),
		UntagMask:                byteMask(prs.Param("untag_mask")),
		VoluntaryCtxtSwitches:    prs.Param("voluntary_ctxt_switches").Uint64(),

		X86ThreadFeatures:       prs.Param("x86_Thread_features").StringList(" "),
		X86ThreadFeaturesLocked: prs.Param("x86_Thread_features_locked").StringList(" "),
	}

	return self, nil
}

func byteMask(p []byte) uint64 {
	p = bytes.TrimPrefix(p, []byte("0x"))
	p = bytes.TrimPrefix(p, []byte("0X"))
	v, _ := strconv.ParseUint(string(p), 16, 64)
	return v
}

func parseHexToBytes(p []byte) []byte {
	s := string(p)
	if len(s)%2 != 0 {
		s = "0" + s
	}

	data, err := hex.DecodeString(s)
	if err != nil {
		return []byte{}
	}
	return data
}

func gidList(p []byte, delim byte) []uint16 {
	pL := bytes.Split(p, []byte{delim})
	data := make([]uint16, len(pL))
	for i, pl := range pL {
		v, _ := strconv.ParseUint(string(pl), 10, 16)
		data[i] = uint16(v)
	}
	return data
}

func parseKbSize(p []byte) uint64 {
	p = bytes.TrimSuffix(p, []byte{' ', 'k', 'B'})
	v, err := strconv.ParseUint(string(p), 10, 64)
	if err != nil {
		return 18446744073709551615
	}
	return v * 1024
}
