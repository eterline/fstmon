package main

import (
	"context"
	"fmt"
	"time"

	systemmonitor "github.com/eterline/fstmon/internal/adapter/system_monitor"
)

func main() {
	cpu := systemmonitor.NewHardwareMetricCPU(2 * time.Second)

	{
		data, err := cpu.CpuPackage()
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
	}

	{
		data, err := cpu.CpuMetrics(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
	}
}
