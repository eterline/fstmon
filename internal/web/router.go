// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package web

import (
	"context"
	"net/http"
	"time"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/services/hostfetchers"
	"github.com/eterline/fstmon/internal/services/monitors"
	"github.com/eterline/fstmon/internal/services/secure"
	"github.com/eterline/fstmon/internal/web/controller"
	"github.com/eterline/fstmon/internal/web/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRouter(ctx context.Context, cfg config.Configuration) http.Handler {

	root := chi.NewMux()

	ext := secure.NewIpExtractor(cfg.ParseIpHeader)

	root.Use(
		middleware.RequestWrapper(ext),
		middleware.RequestLogger,
		middleware.NoCacheControl,
		middleware.SecureControl,
		middleware.SourceSubnetsAllow(ctx, ext, cfg.AllowedSubnets),
		middleware.AllowedHosts(cfg.AllowedHosts),
		middleware.BearerCheck(cfg.AuthToken, ext),
	)

	root.NotFound(controller.NotFound)
	root.MethodNotAllowed(controller.BadMethod)

	root.Route(
		"/monitoring", func(r chi.Router) {

			cpuMon := monitors.InitCpuLoadMon(ctx, 5*time.Second)

			hc := controller.NewHostController(
				hostfetchers.InitSystemMon(ctx, cpuMon),
				hostfetchers.InitAverageLoadMon(ctx),
				hostfetchers.InitPartUseMon(ctx),
				hostfetchers.InitNetworkMon(ctx),
				hostfetchers.InitCpuFetch(cpuMon),
			)

			r.Get("/system", hc.HandleSystem)
			r.Get("/avgload", hc.HandleAvgload)
			r.Get("/parts", hc.HandleParts)
			r.Get("/net", hc.HandleNetworking)
			r.Get("/cpu", hc.HandleCpu)
		},
	)

	return root
}
