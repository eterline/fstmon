package controller

import (
	"log/slog"
	"net/http"

	"github.com/eterline/fstmon/internal/domain"
)

type (
	Fetcher[T any] interface {
		Fetch() (T, error)
	}
)

type HostController struct {
	system     Fetcher[domain.SystemData]
	avgload    Fetcher[domain.AverageLoad]
	partitions Fetcher[domain.PartsUsages]
	network    Fetcher[domain.InterfacesData]
	cpu        Fetcher[domain.CpuLoad]
}

func NewHostController(
	s Fetcher[domain.SystemData],
	a Fetcher[domain.AverageLoad],
	p Fetcher[domain.PartsUsages],
	n Fetcher[domain.InterfacesData],
	c Fetcher[domain.CpuLoad],
) *HostController {
	return &HostController{
		system:     s,
		avgload:    a,
		partitions: p,
		network:    n,
		cpu:        c,
	}
}

func (hc *HostController) HandleNetworking(w http.ResponseWriter, r *http.Request) {

	data, err := hc.network.Fetch()
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

func (hc *HostController) HandleSystem(w http.ResponseWriter, r *http.Request) {

	data, err := hc.system.Fetch()
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

	data, err := hc.partitions.Fetch()
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

	data, err := hc.avgload.Fetch()
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

func (hc *HostController) HandleCpu(w http.ResponseWriter, r *http.Request) {

	data, err := hc.cpu.Fetch()
	if err != nil {
		ResponseError(
			w, http.StatusNotImplemented,
			"could not fetch cpu loads",
		)
		slog.ErrorContext(r.Context(), err.Error())
		return
	}

	ResponseOK(w, data)
}
