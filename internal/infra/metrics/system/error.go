package system

import "fmt"

type SystemError struct {
	embed string
	err   error
}

func (se *SystemError) Error() string {
	if se.err != nil {
		return fmt.Sprintf("system metrics: %s", se.embed)
	}
	return fmt.Sprintf("system metrics: %s: %v", se.embed, se.err)
}

func (se *SystemError) Wrap(err error) *SystemError {
	se.err = err
	return se
}

func (se *SystemError) Unwrap() error {
	return se.err
}

func newSystemError(embed string) *SystemError {
	return &SystemError{
		embed: embed,
	}
}

var (
	ErrScrapeCpuPackage     = newSystemError("failed scrape cpu package")
	ErrScrapeCpuMetrics     = newSystemError("failed scrape cpu metrics")
	ErrScrapeInterfacesIO   = newSystemError("failed scrape interfaces I/O")
	ErrScrapeSystemInfo     = newSystemError("failed scrape system info")
	ErrScrapeMemoryMetrics  = newSystemError("failed scrape memory metrics")
	ErrScrapeThermalMetrics = newSystemError("failed scrape thermal metrics")
	ErrScrapePartitionsInfo = newSystemError("failed scrape partitions info")
	ErrScrapePartitionsIO   = newSystemError("failed scrape partitions I/O")
)
