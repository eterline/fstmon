package system

import (
	"strconv"
	"strings"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/pkg/procf"
)

type hardwareMetricSystem struct{}

func NewHardwareMetricSystem() *hardwareMetricSystem {
	return &hardwareMetricSystem{}
}

func parseProcs(s string) (total, run int) {
	seg := strings.Split(s, "/")
	if len(seg) != 2 {
		return 0, 0
	}

	total, _ = strconv.Atoi(seg[0])
	run, _ = strconv.Atoi(seg[1])

	return
}

func (hms *hardwareMetricSystem) ScrapeSystemInfo() (domain.SystemInfo, error) {
	var data domain.SystemInfo

	up, _ := procf.FetchProcUptime()
	data.Uptime = up.Uptime
	data.Idle = up.IdleTime

	avg, err := procf.FetchProcLoadAvg()
	if err == nil {
		data.AvgLoad.Load1 = avg.Load1
		data.AvgLoad.Load5 = avg.Load5
		data.AvgLoad.Load15 = avg.Load15

		t, r := parseProcs(avg.RunningProcs)
		data.AvgLoad.TotalProcs = t
		data.AvgLoad.RunningProcs = r
	}

	return data, err
}
