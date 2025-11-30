package httphomepage

import (
	"fmt"

	"github.com/eterline/fstmon/internal/domain"
)

// ============================ CPU dto ============================

// DTOCpuCore - simplified per-core dynamic metrics.
type DTOCpuCore struct {
	Load      string `json:"load"`      // "12.5%"
	Frequency string `json:"frequency"` // "3200MHz"
}

// DTOCpu - aggregated CPU info for homepage.
type DTOCpu struct {
	Vendor      string `json:"vendor"`       // "Intel", "AMD"
	Model       string `json:"model"`        // "Ryzen 5 5600X"
	CoreCount   int    `json:"core_count"`   // "6"
	ThreadCount int    `json:"thread_count"` // "12"

	Load      string       `json:"load"`      // "15.4%"
	Frequency string       `json:"frequency"` // "3250MHz"`
	Cores     []DTOCpuCore `json:"cores"`     // per-core metrics
}

func Domain2DTOCpu(pkg domain.CpuPackage, m domain.CpuMetrics) *DTOCpu {
	if len(pkg.Cores) == 0 {
		return &DTOCpu{}
	}

	dto := DTOCpu{
		Vendor:      pkg.Vendor,
		Model:       pkg.ModelName,
		CoreCount:   len(pkg.Cores),
		ThreadCount: pkg.Cores[0].Siblings,
		Load:        fmt.Sprintf("%.1f%%", m.Average.Load),
		Frequency:   fmt.Sprintf("%.1fMhz", m.Average.Frequency),
		Cores:       make([]DTOCpuCore, len(m.Cores)),
	}

	for i, core := range m.Cores {
		dto.Cores[i] = DTOCpuCore{
			Load:      fmt.Sprintf("%.1f%%", core.Load),
			Frequency: fmt.Sprintf("%.1fMhz", core.Frequency),
		}
	}

	return &dto
}

// ============================ RAM dto ============================

type DTOMemoryT struct {
	UsedPercent   string `json:"used_percent"`    // "52.5%"
	Total         string `json:"total"`           // "15.62GB"
	Used          string `json:"used"`            // "8.22GB"
	Free          string `json:"free"`            // "7.32GB"
	UsedTotal     string `json:"used_total"`      // "8.22GB/15.62GB"
	UsedTotalFree string `json:"used_total_free"` // "8.22GB/15.62GB/7.32GB"
}

type DTOMemory struct {
	RAM  DTOMemoryT `json:"ram"`
	Swap DTOMemoryT `json:"swap"`
}

func Domain2DTOMemory(v domain.MemoryMetrics) *DTOMemory {
	return &DTOMemory{
		RAM: DTOMemoryT{
			Total:         NewQBBSBuilder(0).Add(v.Total).Build(),
			Used:          NewQBBSBuilder(0).Add(v.Used).Build(),
			Free:          NewQBBSBuilder(0).Add(v.Free).Build(),
			UsedPercent:   fmt.Sprintf("%.1f%%", v.UsedPercent),
			UsedTotal:     NewQBBSBuilder('/').Add(v.Used).Add(v.Total).Build(),
			UsedTotalFree: NewQBBSBuilder('/').Add(v.Used).Add(v.Total).Add(v.Free).Build(),
		},
		Swap: DTOMemoryT{
			Total:         NewQBBSBuilder(0).Add(v.SwapTotal).Build(),
			Used:          NewQBBSBuilder(0).Add(v.SwapUsed).Build(),
			Free:          NewQBBSBuilder(0).Add(v.SwapFree).Build(),
			UsedPercent:   fmt.Sprintf("%.1f%%", v.SwapUsedPercent),
			UsedTotal:     NewQBBSBuilder('/').Add(v.SwapUsed).Add(v.SwapTotal).Build(),
			UsedTotalFree: NewQBBSBuilder('/').Add(v.SwapUsed).Add(v.SwapTotal).Add(v.SwapFree).Build(),
		},
	}
}

// ============================ System dto ============================

type DTOSystem struct {
	Uptime       string `json:"uptime"`        // "3h-22m-5s"
	Idle         string `json:"idle"`          // "3h-22m-5s"
	Load1        string `json:"load1"`         // "0.42"
	Load5        string `json:"load5"`         // "0.38"
	Load15       string `json:"load15"`        // "0.32"
	RunningProcs int    `json:"running_procs"` // "3"
	TotalProcs   int    `json:"total_procs"`   // "187"
	SummaryProcs string `json:"summary_procs"` // "3/187"
}

func Domain2DTOSystem(v domain.SystemInfo) *DTOSystem {
	return &DTOSystem{
		Uptime:       formatDuration(v.Uptime, ":", true),
		Idle:         formatDuration(v.Idle, ":", true),
		Load1:        fmt.Sprintf("%.2f", v.Load1),
		Load5:        fmt.Sprintf("%.2f", v.Load5),
		Load15:       fmt.Sprintf("%.2f", v.Load15),
		RunningProcs: v.RunningProcs,
		TotalProcs:   v.TotalProcs,
		SummaryProcs: fmt.Sprintf("%d/%d", v.RunningProcs, v.TotalProcs),
	}
}

// ============================ Network dto ============================

type DTONetworkInterfaceIO struct {
	Name string `json:"name"` // "eth0"

	Bytes   IO[uint64] `json:"bytes"`   // "12.3GB"
	Packets IO[uint64] `json:"packets"` // "1.2M"

	BytesSpeed   IO[uint64] `json:"bytes_speed"`   // "120KB/s"
	PacketsSpeed IO[uint64] `json:"packets_speed"` // "133 pkt/s"

	Errors uint64 `json:"errors"` // "0"
	Drops  uint64 `json:"drops"`  // "0"
}

type DTONetworkIO struct {
	BytesTotal   IO[uint64]                       `json:"bytes_total"`
	PacketsTotal IO[uint64]                       `json:"packets_total"`
	ErrorsTotal  uint64                           `json:"errors_total"`
	DropsTotal   uint64                           `json:"drops_total"`
	Interfaces   map[string]DTONetworkInterfaceIO `json:"interfaces"`
}

func NewDTONetworkIO(ifCount int) *DTONetworkIO {
	return &DTONetworkIO{
		Interfaces: make(map[string]DTONetworkInterfaceIO, ifCount),
	}
}

func Domain2DTONetworkInterfaceIO(v domain.InterfacesIOMap) *DTONetworkIO {
	if len(v) == 0 {
		return NewDTONetworkIO(0)
	}

	dto := NewDTONetworkIO(len(v))

	var (
		pktTotalRx   = uint64(0)
		pktTotalTx   = uint64(0)
		bytesTotalRx = uint64(0)
		bytesTotalTx = uint64(0)
		errorsTotal  = uint64(0)
		dropTotal    = uint64(0)
	)

	for name, io := range v {
		pktTotalRx += io.PacketsTotal.RX
		pktTotalTx += io.PacketsTotal.TX
		bytesTotalRx += io.BytesTotal.RX
		bytesTotalTx += io.BytesTotal.TX
		errorsTotal += io.ErrPacketsTotal.Summary
		dropTotal += io.DropPacketsTotal.Summary

		iface := DTONetworkInterfaceIO{
			Name: name,

			Bytes:   NewIOBuilder(io.BytesTotal.RX, io.BytesTotal.TX).AutoUnits().Build(),
			Packets: NewIOBuilder(io.PacketsTotal.RX, io.PacketsTotal.TX).AutoMetricUnits().WithPostfix("Pkts").Build(),

			BytesSpeed:   NewIOBuilder(io.BytesPerSec.RX, io.BytesPerSec.TX).AutoUnitsPerSec().Build(),
			PacketsSpeed: NewIOBuilder(io.PacketsPerSec.RX, io.PacketsPerSec.TX).AutoMetricUnits().WithPostfix("Pkts/s").Build(),

			Errors: io.ErrPacketsTotal.Summary,
			Drops:  io.DropPacketsTotal.Summary,
		}

		dto.Interfaces[name] = iface
	}

	dto.BytesTotal = NewIOBuilder(bytesTotalRx, bytesTotalTx).AutoUnits().Build()
	dto.PacketsTotal = NewIOBuilder(pktTotalRx, pktTotalTx).AutoMetricUnits().WithPostfix("Pkts").Build()
	dto.ErrorsTotal = errorsTotal
	dto.DropsTotal = dropTotal

	return dto
}
