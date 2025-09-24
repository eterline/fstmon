// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

type ProcVmStat map[string]uint64

func ReadProcVmstat() (ProcVmStat, error) {
	data, err := procVmStat.Data()
	if err != nil {
		return nil, err
	}

	prs := NewFileDataSetParser(data, []byte{'\n'}, []byte{' '}, 0, 1)
	res := make(ProcVmStat, prs.Count())

	for key, value := range prs.Data() {
		res[key] = value.Uint64()
	}

	return res, nil
}
