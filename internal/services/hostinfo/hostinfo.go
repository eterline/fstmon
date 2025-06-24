package hostinfo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/pkg/procf"
	pscpu "github.com/shirou/gopsutil/v4/cpu"
	psmem "github.com/shirou/gopsutil/v4/mem"
	psnet "github.com/shirou/gopsutil/v4/net"
)

type HostInfo struct {
	measureIface string
}

func InitHostInfo(iface string) *HostInfo {
	return &HostInfo{
		measureIface: iface,
	}
}

func (hi *HostInfo) System(ctx context.Context) (domain.SystemData, error) {

	loads, err := pscpu.PercentWithContext(ctx, 5*time.Second, true)
	if err != nil {
		return domain.SystemData{}, err
	}

	stat, err := psmem.VirtualMemory()
	if err != nil {
		return domain.SystemData{}, err
	}

	uptime, err := procf.FetchProcUptime()
	if err != nil {
		return domain.SystemData{}, err
	}

	temp, err := cpuTemperature(ctx)
	if err != nil {
		return domain.SystemData{}, err
	}

	data := domain.SystemData{
		CpuTemp: fmt.Sprintf("%.1f°C", temp),
		Uptime:  formatTime(uptime.Uptime),
		Cpu:     fmt.Sprintf("%.0f%%", avgFloat64(loads)),
		Memory:  fmt.Sprintf("%.0f%%", stat.UsedPercent),
	}

	return data, nil
}

func (hi *HostInfo) Processes(ctx context.Context) (domain.ProcessesData, error) {
	return domain.ProcessesData{}, nil
}

func (hi *HostInfo) Networking(ctx context.Context) (domain.NetworkingData, error) {
	data := domain.NetworkingData{}

	err := calcRxTx(ctx, hi.measureIface, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func cpuTemperature(ctx context.Context) (float64, error) {

	data, err := procf.Temperatures(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	temp := float64(0)
	contains := func(value string) bool {
		return strings.Contains(value, "coretemp") ||
			strings.Contains(value, "k10temp")
	}

	for key, value := range data {
		if !contains(key) {
			continue
		}
		count++
		temp += value.Current
	}

	if count == 0 {
		return 0, nil
	}

	return (temp / float64(count)), nil
}

func calcRxTx(ctx context.Context, iface string, data *domain.NetworkingData) error {

	counters1, err := psnet.IOCounters(true)
	if err != nil {
		return err
	}

	var rx1, tx1 uint64
	found := false
	for _, c := range counters1 {
		if c.Name == iface {
			rx1 = c.BytesRecv
			tx1 = c.BytesSent
			found = true
			break
		}
	}
	if !found {
		*data = domain.NetworkingData{
			FullRX:  "null",
			FullTX:  "null",
			SpeedRX: "null",
			SpeedTX: "null",
		}
		return nil
	}

	select {
	case <-ctx.Done():
		return errors.New("measurement cancelled by context")
	case <-time.After(1 * time.Second):

	}

	counters2, err := psnet.IOCounters(true)
	if err != nil {
		return err
	}

	var rx2, tx2 uint64
	found = false
	for _, c := range counters2 {
		if c.Name == iface {
			rx2 = c.BytesRecv
			tx2 = c.BytesSent
			found = true
			break
		}
	}
	if !found {
		return errors.New("interface not found on second read: " + iface)
	}

	data.FullRX = bytesString(int64(rx2))
	data.FullTX = bytesString(int64(tx2))
	data.SpeedRX = bytesStringSpeed(int64(rx2 - rx1))
	data.SpeedTX = bytesStringSpeed(int64(tx2 - tx1))

	return nil
}

func bytesString(v int64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
	)

	fv := float64(v)

	switch {
	case fv >= TB:
		return fmt.Sprintf("%.2fTB", fv/TB)
	case fv >= GB:
		return fmt.Sprintf("%.2fGB", fv/GB)
	case fv >= MB:
		return fmt.Sprintf("%.2fMB", fv/MB)
	case fv >= KB:
		return fmt.Sprintf("%.2fKB", fv/KB)
	default:
		return fmt.Sprintf("%dB", v)
	}
}

func avgFloat64(l []float64) float64 {
	len := len(l)
	if len == 0 {
		return 0.0
	}

	sum := float64(0)
	for _, v := range l {
		sum += v
	}
	return sum / float64(len)
}

func bytesStringSpeed(v int64) string {
	return bytesString(v) + "/s"
}

func formatTime(t time.Duration) string {
	seconds := int(t.Seconds())

	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60

	res := ""

	if h > 0 {
		res += fmt.Sprintf("%dh", h)
	}
	if m > 0 {
		res += fmt.Sprintf("%dm", m)
	}
	if s > 0 || res == "" {
		res += fmt.Sprintf("%ds", s)
	}

	return res
}
