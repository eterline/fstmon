// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package convert

import (
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/interface/grpc/flugel/common"
	"google.golang.org/protobuf/types/known/durationpb"
)

func toIOUint64(io domain.IO[uint64]) *common.IOUint64 {
	return &common.IOUint64{
		Summary: io.Summary,
		Rx:      io.RX,
		Tx:      io.TX,
	}
}

func toIODuration(d domain.IO[time.Duration]) *common.IODuration {
	return &common.IODuration{
		Summary: durationpb.New(d.Summary),
		Rx:      durationpb.New(d.RX),
		Tx:      durationpb.New(d.TX),
	}
}
