// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

import (
	"strconv"
	"testing"
)

func Test_bytesToInt32(t *testing.T) {
	for i := int64(-32000); i < 32000; i++ {
		v := strconv.FormatInt(i, 10)
		if out := bytesToInt32([]byte(v)); out != int32(i) {
			t.Error("bytesToInt32", i, out)
		}
	}
}

func Test_bytesToInt64(t *testing.T) {
	for i := int64(-32000); i < 32000; i++ {
		v := strconv.FormatInt(i, 10)
		if out := bytesToInt64([]byte(v)); out != i {
			t.Error("bytesToInt64", i, out)
		}
	}
}

func Test_bytesToUint32(t *testing.T) {
	for i := uint64(0); i < 32000; i++ {
		v := strconv.FormatUint(i, 10)
		if out := bytesToUint32([]byte(v)); out != uint32(i) {
			t.Error("bytesToUint32", i, out)
		}
	}
}

func Test_bytesToUint64(t *testing.T) {
	for i := uint64(0); i < 32000; i++ {
		v := strconv.FormatUint(i, 10)
		if out := bytesToUint64([]byte(v)); out != i {
			t.Error("bytesToUint64", i, out)
		}
	}
}
