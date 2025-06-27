package controller

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/eterline/fstmon/internal/domain"
)

type HostDataProvider interface {
	Networking(context.Context) (domain.InterfacesData, error)
	Processes(context.Context) (domain.ProcessesData, error)
	System(context.Context) (domain.SystemData, error)
	PartUse(context.Context) (domain.PartsUsages, error)
	AverageLoad() (domain.AverageLoad, error)
	TemperatureMap(context.Context) (domain.TemperatureMap, error)
}

type CpuProvider interface {
	CpuLoad() (domain.CpuLoad, error)
}

type HostController struct {
	hostData HostDataProvider
	cpuData  CpuProvider
}

func NewHostController(hdp HostDataProvider, cp CpuProvider) *HostController {
	return &HostController{
		hostData: hdp,
		cpuData:  cp,
	}
}

func (hc *HostController) HandleNetworking(w http.ResponseWriter, r *http.Request) {

	data, err := hc.hostData.Networking(r.Context())
	if err != nil {
		ResponseError(
			w, http.StatusNotImplemented,
			"could not fetch network counters",
		)
		slog.ErrorContext(r.Context(), err.Error())
		return
	}

	ResponseOK(w, data)
}

func (hc *HostController) HandleSys(w http.ResponseWriter, r *http.Request) {

	data, err := hc.hostData.System(r.Context())
	if err != nil {
		ResponseError(
			w, http.StatusNotImplemented,
			"could not fetch system info",
		)
		slog.ErrorContext(r.Context(), err.Error())
		return
	}

	ResponseOK(w, data)
}

func (hc *HostController) HandleParts(w http.ResponseWriter, r *http.Request) {

	data, err := hc.hostData.PartUse(r.Context())
	if err != nil {
		ResponseError(
			w, http.StatusNotImplemented,
			"could not fetch parts info",
		)
		slog.ErrorContext(r.Context(), err.Error())
		return
	}

	ResponseOK(w, data)
}

func (hc *HostController) HandleAvgload(w http.ResponseWriter, r *http.Request) {

	data, err := hc.hostData.AverageLoad()
	if err != nil {
		ResponseError(
			w, http.StatusNotImplemented,
			"could not fetch avg load",
		)
		slog.ErrorContext(r.Context(), err.Error())
		return
	}

	ResponseOK(w, data)
}

func (hc *HostController) HandleTemp(w http.ResponseWriter, r *http.Request) {

	data, err := hc.hostData.TemperatureMap(r.Context())
	if err != nil {
		ResponseError(
			w, http.StatusNotImplemented,
			"could not fetch temperature sensors data",
		)
		slog.ErrorContext(r.Context(), err.Error())
		return
	}

	ResponseOK(w, data)
}

func (hc *HostController) HandleCpu(w http.ResponseWriter, r *http.Request) {

	data, err := hc.cpuData.CpuLoad()
	if err != nil {
		ResponseError(
			w, http.StatusNotImplemented,
			"could not fetch cpu data",
		)
		slog.ErrorContext(r.Context(), err.Error())
		return
	}

	ResponseOK(w, data)
}
