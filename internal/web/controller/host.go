// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package controller

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/log"
)

type (
	Fetcher[T any] interface {
		Fetch() (T, error)
	}
)

type HostController struct {
	log        *slog.Logger
	system     Fetcher[domain.SystemData]
	avgload    Fetcher[domain.AverageLoad]
	partitions Fetcher[domain.PartsUsages]
	network    Fetcher[domain.InterfacesData]
	cpu        Fetcher[domain.CpuLoad]
}

func NewHostController(
	ctx context.Context,
	s Fetcher[domain.SystemData],
	a Fetcher[domain.AverageLoad],
	p Fetcher[domain.PartsUsages],
	n Fetcher[domain.InterfacesData],
	c Fetcher[domain.CpuLoad],
) *HostController {
	return &HostController{
		log:        log.MustLoggerFromContext(ctx),
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
		hc.log.ErrorContext(r.Context(), err.Error())
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
		hc.log.ErrorContext(r.Context(), err.Error())
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
		hc.log.ErrorContext(r.Context(), err.Error())
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
		hc.log.ErrorContext(r.Context(), err.Error())
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
		hc.log.ErrorContext(r.Context(), err.Error())
		return
	}

	ResponseOK(w, data)
}
