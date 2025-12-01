// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package app

import (
	"context"
	"time"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/infra/http/common/api"
	"github.com/eterline/fstmon/internal/infra/http/common/security"
	httphomepage "github.com/eterline/fstmon/internal/infra/http/homepage"
	middleware "github.com/eterline/fstmon/internal/infra/http/middlewares"
	"github.com/eterline/fstmon/internal/infra/http/server"
	metricstore "github.com/eterline/fstmon/internal/infra/metrics/metric_store"
	"github.com/eterline/fstmon/internal/infra/metrics/system"
	"github.com/eterline/fstmon/internal/log"
	"github.com/eterline/fstmon/internal/services/monitor"
	"github.com/eterline/fstmon/pkg/toolkit"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/procfs"
)

type InitFlags struct {
	CommitHash string
	Version    string
}

func Execute(root *toolkit.AppStarter, flags InitFlags, cfg config.Configuration) {
	ctx := root.Context
	log := log.MustLoggerFromContext(ctx)

	log.Info("starting app", "commit", flags.CommitHash, "version", flags.Version)
	defer func() {
		wt := root.WorkTime()
		log.Info("exit from app", "running_time", wt)
	}()

	// ============================

	proc, err := procfs.NewDefaultFS()
	if err != nil {
		log.Error("procfs initialization error", "error", err)
		root.MustStopApp(1)
	}

	tokenAuth, err := security.NewTokenAuthProvide(security.PolicyMid, cfg.AuthToken...)
	if err != nil {
		log.Error("token auth initialization error", "error", err)
		root.MustStopApp(1)
	}

	log.Warn("token auth setup", "enabled", tokenAuth.Enabled(), "registered", len(cfg.AuthToken))

	// ============================

	log.Info("in-memory pooler store initialization")
	mStore := metricstore.NewMetricInMemoryStore()
	defer func() {
		mStore.Close()
		log.Info("in-memory pooler store closed")
	}()

	metricPooling := monitor.NewServicePooler(mStore) // Metric pooling service

	// ============================

	// CPU usage
	hMtCpu := system.NewHardwareMetricCPU(time.Second)
	metricPooling.AddMetricPooling(
		wrapJob(hMtCpu.ScrapeCpuMetrics), "cpu_usage", cfg.CpuDuration(),
	)

	// CPU summary info
	metricPooling.AddMetricPooling(
		wrapJob(hMtCpu.ScrapeCpuPackage), "cpu", time.Minute,
	)

	// Network interfaces I/O
	hMtNet := system.NewHardwareMetricNetwork(proc)
	metricPooling.AddMetricPooling(
		wrapJob(hMtNet.ScrapeInterfacesIO), "net_io", cfg.NetworkIoDuration(),
	)

	// Memory stats
	hMtMem := system.NewHardwareMetricMemory(proc)
	metricPooling.AddMetricPooling(
		wrapJob(hMtMem.ScrapeMemoryMetrics), "memory", cfg.MemorykDuration(),
	)

	hMtParts := system.NewHardwareMetricPartitions(proc)

	// Disk I/O metrics
	metricPooling.AddMetricPooling(
		wrapJob(hMtParts.ScrapeDiskIO), "disk_io", cfg.DiskIODuration(),
	)

	// Partitions static info
	metricPooling.AddMetricPooling(
		wrapJob(hMtParts.ScrapePartitions), "partitions", time.Minute,
	)

	// System info (loadavg, uptime, procs, etc.)
	hMtSystem := system.NewHardwareMetricSystem(proc)
	metricPooling.AddMetricPooling(
		wrapJob(hMtSystem.ScrapeSystemInfo), "system", cfg.SystemDuration(),
	)

	// Thermal sensors
	hMtThermal := system.NewHardwareThermalMetrics()
	metricPooling.AddMetricPooling(
		wrapJob(hMtThermal.ScrapeThermalMetrics), "thermal", cfg.ThermalDuration(),
	)

	// ============================

	root.WrapWorker(func() {
		metricPooling.RunPooling(ctx)
		metricPooling.Wait()
		log.Info("metric pooling stopped")
	})

	// ============================

	rootMux := chi.NewMux()

	// change default handlers
	rootMux.NotFound(api.HandleNotFound)
	rootMux.MethodNotAllowed(api.HandleNotAllowedMethod)

	rootMux.Use(middleware.NoCacheControl)
	rootMux.Use(middleware.RequestLoggerWrap(log))
	rootMux.Use(middleware.AllowedHosts(ctx, cfg.AllowedHosts))
	rootMux.Use(middleware.BearerAuth(ctx, tokenAuth)) // lookup at NewTokenAuthProvide upper

	// ============================

	metricRouter := chi.NewRouter()
	h := httphomepage.New(metricPooling)

	metricRouter.Route("/homepage",
		func(r chi.Router) {
			r.Get("/network", h.HandleNetwork)
			r.Get("/memory", h.HandleMemory)
			r.Get("/system", h.HandleSystem)
			r.Get("/cpu", h.HandleCpu)
			r.Get("/partitions", h.HandlePartitions)
			r.Get("/diskio", h.HandleDiskIO)
		},
	)

	rootMux.Mount("/metric", metricRouter)

	// ============================

	srv := server.NewServer(
		rootMux,
		server.WithTLS(server.NewServerTlsConfig()),
		server.WithDisabledDefaultHttp2Map(),
	)
	defer srv.Close()

	// ============================

	root.WrapWorker(func() {
		err := srv.Run(ctx, cfg.Listen, cfg.KeyFileSSL, cfg.CrtFileSSL)
		if err != nil {
			log.Error("server run error", "error", err)
		}
	})

	// ============================
	root.WaitWorkers(15 * time.Second)
}

func wrapJob[T any](f func(context.Context) (T, error)) monitor.UpdateWorker {
	return func(ctx context.Context) (any, error) {
		return f(ctx)
	}
}
