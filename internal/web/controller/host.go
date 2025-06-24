package controller

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/eterline/fstmon/internal/domain"
)

type HostDataProvider interface {
	Networking(context.Context) (domain.NetworkingData, error)
	Processes(context.Context) (domain.ProcessesData, error)
	System(context.Context) (domain.SystemData, error)
}

type HostController struct {
	hostData HostDataProvider
}

func NewHostController(hdp HostDataProvider) *HostController {
	return &HostController{
		hostData: hdp,
	}
}

func (hc *HostController) HandleNetworking(w http.ResponseWriter, r *http.Request) {

	data, err := hc.hostData.Networking(r.Context())
	if err != nil {
		ResponseError(
			w, http.StatusNotImplemented,
			"Could not fetch network counters",
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
			"Could not fetch network counters",
		)
		slog.ErrorContext(r.Context(), err.Error())
		return
	}

	ResponseOK(w, data)
}
