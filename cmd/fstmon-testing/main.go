package main

import (
	"context"

	systemmonitor "github.com/eterline/fstmon/internal/infra/metrics/system"
	"github.com/eterline/fstmon/internal/utils/output"
	"github.com/prometheus/procfs"
)

func main() {

	fs, err := procfs.NewDefaultFS()
	if err != nil {
		panic(err)
	}

	net := systemmonitor.NewHardwareMetricPartitions(fs)

	data, err := net.ScrapePartitionsInfo(context.Background())
	if err != nil {
		panic(err)
	}
	output.PrintlnPrettyJSON(data)

}
