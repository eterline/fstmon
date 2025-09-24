package hostfetchers

import (
	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/services/monitors"
)

type CpuFetch struct {
	mon *monitors.CpuLoadMonitoring
}

func InitCpuFetch(mon *monitors.CpuLoadMonitoring) *CpuFetch {
	return &CpuFetch{
		mon: mon,
	}
}

func (c *CpuFetch) Fetch() (domain.CpuLoad, error) {
	return c.mon.Data()
}
