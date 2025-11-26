// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package web

import (
	"context"
	"net/http"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/log"
	"github.com/eterline/fstmon/internal/services/hostfetchers"
	"github.com/eterline/fstmon/internal/services/monitors"
	"github.com/eterline/fstmon/internal/services/secure"
	"github.com/eterline/fstmon/internal/web/controller"
	"github.com/eterline/fstmon/internal/web/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRouter(ctx context.Context, cfg config.Configuration) http.Handler {
	log := log.MustLoggerFromContext(ctx)

	log.Info(
		"starting monitoring loop",
		"cpu", cfg.CpuDuration(),
		"system", cfg.SystemDuration(),
		"avgload", cfg.AvgloadDuration(),
		"partitions", cfg.PartitionsDuration(),
		"network", cfg.NetworkDuration(),
	)

	cpuMon := monitors.InitCpuLoadMon(ctx, cfg.CpuDuration())

	hc := controller.NewHostController(
		hostfetchers.InitSystemMon(ctx, cpuMon, cfg.SystemDuration()),
		hostfetchers.InitAverageLoadMon(ctx, cfg.AvgloadDuration()),
		hostfetchers.InitPartUseMon(ctx, cfg.PartitionsDuration()),
		hostfetchers.InitNetworkMon(ctx, cfg.NetworkDuration()),
		hostfetchers.InitCpuFetch(cpuMon),
	)

	ext := secure.NewIpExtractor(cfg.ParseIpHeader)

	root := chi.NewMux()

	root.Use(
		middleware.RequestWrapper(ext),
		middleware.RequestLogger(log),
		middleware.NoCacheControl,
		middleware.SecureControl,
		middleware.SourceSubnetsAllow(ctx, ext, cfg.AllowedSubnets),
		middleware.AllowedHosts(cfg.AllowedHosts),
		middleware.BearerCheck(ctx, cfg.AuthToken, ext),
	)

	root.NotFound(controller.NotFound)
	root.MethodNotAllowed(controller.BadMethod)

	root.Route(
		"/monitoring", func(r chi.Router) {

			r.Get("/system", hc.HandleSystem)
			r.Get("/avgload", hc.HandleAvgload)
			r.Get("/parts", hc.HandleParts)
			r.Get("/net", hc.HandleNetworking)
			r.Get("/cpu", hc.HandleCpu)
		},
	)

	return root
}
