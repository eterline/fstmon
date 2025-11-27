package main

import (
	"context"
	"time"

	systemmonitor "github.com/eterline/fstmon/internal/infra/monitoring/system"
	"github.com/eterline/fstmon/internal/utils/output"
)

func main() {
	net := systemmonitor.NewHardwareMetricNetwork()

	t := time.NewTicker(2 * time.Second)
	defer t.Stop()

	for range t.C {
		data, err := net.ScrapeInterfacesIO(context.Background())
		if err != nil {
			panic(err)
		}
		output.PrintlnPrettyJSON(data)
	}

}
