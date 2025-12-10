// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
SizesTLB – handle Translation Lookaside Buffer

	┌───────────┬────────────────────────────────────────────────────┐
	│ Field     │ Description                                        │
	├───────────┼────────────────────────────────────────────────────┤
	│ Writes    │ Number of TLB write operations or entries          │
	│ PageSize  │ Page size in bytes associated with this TLB entry  │
	└───────────┴────────────────────────────────────────────────────┘
*/
type SizesTLB struct {
	Writes   uint64 `json:"writes"`
	PageSize uint64 `json:"page_size"`
}

/*
AddrBits – handle memory address bits

	┌──────────┬────────────────────────────────────────────────────────────┐
	│ Field    │ Description                                                │
	├──────────┼────────────────────────────────────────────────────────────┤
	│ Physical │ Number of bits used for physical memory addressing         │
	│ Virtual  │ Number of bits used for virtual memory addressing          │
	└──────────┴────────────────────────────────────────────────────────────┘
*/
type AddrBits struct {
	Physical uint8 `json:"physical"`
	Virtual  uint8 `json:"virtual"`
}

type Instructions []string

func (is Instructions) Count() int {
	return len(is)
}

func (is Instructions) HaveInstruction(ins string) bool {
	ins = strings.ToLower(ins)
	l := is.Count()

	for i := 0; i < l; i++ {
		if is[i] == ins {
			return true
		}
	}

	return false
}

func FetchCpuInfo() (ProcCpuInfo, error) {

	info := ProcCpuInfo{
		Cores: []CoreInfo{},
	}

	data, err := procCpuInfo.Data()
	if err != nil {
		return info, fmt.Errorf("cpu info fetch error: %v", err)
	}

	var (
		delimSetsBytes  = []byte{'\n', '\n'}
		delimSetBytes   = []byte{'\n'}
		delimParamBytes = []byte{':', ' '}
	)

	dataSets := bytes.Split(data, delimSetsBytes)
	setsCount := len(dataSets) - 1
	if setsCount < 1 {
		return info, fmt.Errorf("invalid info file: '%s'", procCpuInfo)
	}

	dataSets = dataSets[:setsCount]
	info.Cores = make([]CoreInfo, setsCount)

	for i, set := range dataSets {

		prs := NewFileDataSetParser(set, delimSetBytes, delimParamBytes, 0, 1)

		core := CoreInfo{
			Processor:       prs.Param("processor").Int64(),
			VendorID:        prs.Param("vendor_id").String(),
			CpuFamily:       prs.Param("cpufamily").Int64(),
			Model:           prs.Param("model").Int64(),
			ModelName:       prs.Param("modelname").String(),
			Stepping:        prs.Param("stepping").Int64(),
			CpuMhz:          prs.Param("cpuMHz").Float64(),
			PhysicalID:      prs.Param("physicalid").Int64(),
			Siblings:        prs.Param("siblings").Int64(),
			CoreID:          prs.Param("coreid").Int64(),
			CpuCores:        prs.Param("cpucores").Int64(),
			Apicid:          prs.Param("apicid").Int64(),
			InitialApicid:   prs.Param("initialapicid").Int64(),
			CpuIDLevel:      prs.Param("cpuidlevel").Int64(),
			Flags:           prs.Param("flags").StringList(" "),
			Bugs:            prs.Param("bugs").StringList(" "),
			BogoMips:        prs.Param("bogomips").Float64(),
			ClflushSize:     prs.Param("clflushsize").Int64(),
			CacheAlignment:  prs.Param("cache_alignment").Int64(),
			FpuException:    prs.Param("fpu_exception").ParseBool("yes"),
			Fpu:             prs.Param("fpu").ParseBool("yes"),
			WP:              prs.Param("wp").ParseBool("yes"),
			PowerManagement: prs.Param("powermanagement").StringList(" "),
			VmxFlags:        prs.Param("vmxflags").StringList(" "),
			Microcode:       prs.Param("microcode").String(),
			CacheSize:       cpuCache(prs.Param("cachesize")),
			TLB:             cpuTLB(prs.Param("TLBsize")),
			AddressSizes:    cpuMemAddrs(prs.Param("addresssizes")),
		}

		info.Cores[i] = core
	}

	return info, nil
}

func cpuMemAddrs(p []byte) AddrBits {
	var data AddrBits

	re := regexp.MustCompile(`(\d+)\s+bits\s+physical.*?(\d+)\s+bits\s+virtual`)
	matches := re.FindStringSubmatch(string(p))
	if len(matches) < 3 {
		return data
	}

	data.Physical = uint8(bytesToInt32([]byte(matches[1])))
	data.Virtual = uint8(bytesToInt32([]byte(matches[2])))

	return data
}

func cpuTLB(p []byte) SizesTLB {
	fields := bytes.Fields(p)
	if len(fields) < 2 {
		return SizesTLB{}
	}

	writes, err := strconv.ParseUint(string(fields[0]), 10, 64)
	if err != nil {
		return SizesTLB{}
	}

	pageSizeStr := strings.ToUpper(string(fields[1]))
	var multiplier uint64

	switch {
	case strings.HasSuffix(pageSizeStr, "K"):
		multiplier = 1024
	case strings.HasSuffix(pageSizeStr, "M"):
		multiplier = 1024 * 1024
	case strings.HasSuffix(pageSizeStr, "G"):
		multiplier = 1024 * 1024 * 1024
	default:
		return SizesTLB{}
	}

	pageNumStr := strings.TrimRight(pageSizeStr, "KMG")
	pageNum, err := strconv.ParseUint(pageNumStr, 10, 64)
	if err != nil {
		return SizesTLB{}
	}

	return SizesTLB{
		Writes:   writes,
		PageSize: pageNum * multiplier,
	}
}

func cpuCache(p []byte) int64 {
	p = bytes.Fields(p)[0]
	kbCount := bytesToInt64(p)
	return kbCount * 1024
}

func cpuMicrocode(p []byte) []byte {
	st := string(p)
	data, _ := hex.DecodeString(st)
	return data
}
