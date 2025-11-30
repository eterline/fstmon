// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package app

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/infra/http/common/api"
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

	mStore := metricstore.NewMetricInMemoryStore()
	defer mStore.Close()

	metricPooling := monitor.NewServicePooler(mStore)

	// ===========

	hMtCpu := system.NewHardwareMetricCPU(time.Second)
	metricPooling.AddMetricPooling(
		WrapJob(hMtCpu.ScrapeCpuMetrics), "cpu_dynamic", cfg.CpuDuration(),
	)

	hMtNet := system.NewHardwareMetricNetwork(proc)
	metricPooling.AddMetricPooling(
		WrapJob(hMtNet.ScrapeInterfacesIO), "net_io", cfg.NetworkDuration(),
	)

	hMtMem := system.NewHardwareMetricMemory(proc)
	metricPooling.AddMetricPooling(
		WrapJob(hMtMem.ScrapeMemoryMetrics), "memory", cfg.MemorykDuration(),
	)

	hMtParts := system.NewHardwareMetricPartitions(proc)
	metricPooling.AddMetricPooling(
		WrapJob(hMtParts.ScrapePartitionsInfo), "partitions", cfg.PartitionsDuration(),
	)

	hMtSystem := system.NewHardwareMetricSystem(proc)
	metricPooling.AddMetricPooling(
		WrapJob(hMtSystem.ScrapeSystemInfo), "system", cfg.SystemDuration(),
	)

	hMtThermal := system.NewHardwareThermalMetrics()
	metricPooling.AddMetricPooling(
		WrapJob(hMtThermal.ScrapeThermalMetrics), "thermal", cfg.ThermalDuration(),
	)

	// ===========

	metricPooling.RunPooling(ctx)

	root.NewThread()
	go func() {
		defer root.DoneThread()
		metricPooling.Wait()
		log.Info("metric pooling stopped")
	}()

	// ============================

	// TODO: http adapter setup

	// ============================

	m := chi.NewMux()

	m.Get("/metric/{metricKey}",
		func(w http.ResponseWriter, r *http.Request) {
			metricKey := chi.URLParam(r, "metricKey")
			m, wkExists, mtExists, retryIn := metricPooling.ActualMetric(metricKey)

			if !wkExists {
				api.NewResponse().
					SetCode(http.StatusNotFound).
					SetMessage("uncorrect metric key").
					AddStringError("metric worker not exists").
					Write(w)

				return
			}

			if !mtExists {
				w.Header().Set("Retry-After", strconv.Itoa(int(retryIn.Seconds())))

				api.NewResponse().
					SetCode(http.StatusServiceUnavailable).
					SetMessage("metric not exists").
					AddStringError("metric did not scraped yet").
					Write(w)

				return
			}

			api.NewResponse().SetCode(http.StatusOK).WrapData(m).Write(w)
		},
	)

	srv := server.NewServer(m)

	// ============================

	root.NewThread()
	go func() {
		defer root.DoneThread()

		err := srv.Run(ctx, ":3000", "", "")
		if err != nil {
			log.Error("server run error", "error", err)
		}
	}()

	// ============================
	root.Wait()
	root.WaitThreads(15 * time.Second)
	srv.Close()
}

func WrapJob[T any](f func(context.Context) (T, error)) monitor.UpdateWorker {
	return func(ctx context.Context) (any, error) {
		return f(ctx)
	}
}
