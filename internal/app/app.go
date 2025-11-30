// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package app

import (
	"context"
	"time"

	"github.com/eterline/fstmon/internal/config"
	httphomepage "github.com/eterline/fstmon/internal/infra/http/homepage"
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
		log.Error("procfs init error", "error", err)
		root.MustStopApp(1)
	}

	// tokenAuth, err := security.NewTokenAuthProvide(security.PolicyMid, cfg.AuthToken)
	// if err != nil {
	// 	log.Error("token auth initialization error", "error", err)
	// 	root.MustStopApp(1)
	// }

	// ============================

	log.Info("in-memory store init")
	mStore := metricstore.NewMetricInMemoryStore()
	defer func() {
		log.Info("in-memory store closed")
		mStore.Close()
	}()

	metricPooling := monitor.NewServicePooler(mStore) // Metric pooling service

	// ===========

	// CPU usage
	hMtCpu := system.NewHardwareMetricCPU(time.Second)
	metricPooling.AddMetricPooling(
		WrapJob(hMtCpu.ScrapeCpuMetrics), "cpu_usage", cfg.CpuDuration(),
	)

	// CPU summary info
	metricPooling.AddMetricPooling(
		WrapJob(hMtCpu.ScrapeCpuMetrics), "cpu", time.Minute,
	)

	// Network interfaces I/O
	hMtNet := system.NewHardwareMetricNetwork(proc)
	metricPooling.AddMetricPooling(
		WrapJob(hMtNet.ScrapeInterfacesIO), "net_io", cfg.NetworkIoDuration(),
	)

	// Memory stats
	hMtMem := system.NewHardwareMetricMemory(proc)
	metricPooling.AddMetricPooling(
		WrapJob(hMtMem.ScrapeMemoryMetrics), "memory", cfg.MemorykDuration(),
	)

	// Partitions static info
	hMtParts := system.NewHardwareMetricPartitions(proc)
	metricPooling.AddMetricPooling(
		WrapJob(hMtParts.ScrapePartitionsInfo), "partitions", time.Minute,
	)

	// Partitions I/O metrics
	metricPooling.AddMetricPooling(
		WrapJob(hMtParts.ScrapePartitionIO), "partitions_io", cfg.PartitionsIoDuration(),
	)

	// System info (loadavg, uptime, procs, etc.)
	hMtSystem := system.NewHardwareMetricSystem(proc)
	metricPooling.AddMetricPooling(
		WrapJob(hMtSystem.ScrapeSystemInfo), "system", cfg.SystemDuration(),
	)

	// Thermal sensors
	hMtThermal := system.NewHardwareThermalMetrics()
	metricPooling.AddMetricPooling(
		WrapJob(hMtThermal.ScrapeThermalMetrics), "thermal", cfg.ThermalDuration(),
	)

	// ===========

	root.WrapWorker(func() {
		metricPooling.RunPooling(ctx)
		metricPooling.Wait()
		log.Info("metric pooling stopped")
	})

	// ============================

	rootMux := chi.NewMux()

	metricRouter := chi.NewRouter()

	// =========

	// TODO:
	h := httphomepage.New(metricPooling, log)
	metricRouter.Route("/homepage",
		func(r chi.Router) {
			r.Get("/sensors", h.HandleThermal)
			r.Get("/system", h.HandleSystem)
		},
	)

	// =========

	metricRouter.Route("/flugel",
		func(r chi.Router) {
			// TODO: wiil be later
		},
	)

	// =========

	rootMux.Mount("/metric", metricRouter)

	// ============================

	srv := server.NewServer(rootMux, server.WithTLS(server.NewServerTlsConfig()))
	defer srv.Close()

	// ============================

	root.WrapWorker(func() {
		err := srv.Run(ctx, ":3000", cfg.KeyFileSSL, cfg.CrtFileSSL)
		if err != nil {
			log.Error("server run error", "error", err)
		}
	})

	// ============================
	root.WaitWorkers(15 * time.Second)
}

func WrapJob[T any](f func(context.Context) (T, error)) monitor.UpdateWorker {
	return func(ctx context.Context) (any, error) {
		return f(ctx)
	}
}
