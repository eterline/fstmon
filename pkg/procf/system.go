// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

import (
	"bytes"
	"fmt"
	"time"
)

func FetchProcLoadAvg() (ProcLoadAvg, error) {
	data, err := procLoadAvg.Data()
	if err != nil {
		return ProcLoadAvg{}, err
	}

	dataFlds := bytes.Fields(data)
	if len(dataFlds) < 5 {
		return ProcLoadAvg{}, fmt.Errorf("invalid proc file")
	}

	stat := ProcLoadAvg{
		Load1:        bytesToFloat64(dataFlds[0]),
		Load5:        bytesToFloat64(dataFlds[1]),
		Load15:       bytesToFloat64(dataFlds[2]),
		RunningProcs: bytesToString(dataFlds[3]),
		LastPID:      bytesToInt32(dataFlds[4]),
	}

	return stat, nil
}

func FetchProcUptime() (ProcUptime, error) {
	data, err := procUptime.Data()
	if err != nil {
		return ProcUptime{}, err
	}

	dataFlds := bytes.Fields(data)
	if len(dataFlds) < 2 {
		return ProcUptime{}, fmt.Errorf("invalid proc file")
	}

	stat := ProcUptime{
		Uptime:   time.Duration(bytesToFloat64(dataFlds[0])*1000.0) * time.Millisecond,
		IdleTime: time.Duration(bytesToFloat64(dataFlds[1])*1000.0) * time.Millisecond,
	}

	return stat, nil
}

func FetchProcCrypto() (ProcCrypto, error) {
	data, err := procCrypto.Data()
	if err != nil {
		return ProcCrypto{}, err
	}

	dataSets := bytes.Split(data, []byte{'\n', '\n'})
	if len(dataSets) < 1 {
		return ProcCrypto{}, fmt.Errorf("invalid proc file")
	}
	dataSets = dataSets[:len(dataSets)-1]
	res := make(ProcCrypto, len(dataSets))

	for i, set := range dataSets {
		prs := NewFileDataSetParser(set, []byte{'\n'}, []byte{':', ' '}, 0, 1)

		res[i] = ProcCryptoModule{
			Name:     prs.Param("name").String(),
			Driver:   prs.Param("driver").String(),
			Module:   prs.Param("module").String(),
			Priority: prs.Param("priority").Int32(),
			Refcnt:   prs.Param("refcnt").Int32(),
			SelfTest: prs.Param("selftest").String(),
			Internal: prs.Param("internal").String(),
			Type:     prs.Param("type").String(),
		}
	}

	return res, nil
}
