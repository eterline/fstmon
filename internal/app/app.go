// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package app

import (
	"net/http"
	"time"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/infra/http/common/api"
	"github.com/eterline/fstmon/internal/infra/http/common/security"
	middleware "github.com/eterline/fstmon/internal/infra/http/middlewares"
	"github.com/eterline/fstmon/internal/infra/http/server"
	"github.com/eterline/fstmon/internal/infra/metrics/system"
	"github.com/eterline/fstmon/internal/log"
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

	pfs, err := procfs.NewDefaultFS()
	if err != nil {
		log.Error("procfs init error", "error", err)
	}
	net := system.NewHardwareMetricNetwork(pfs)

	m := chi.NewMux()
	m.Use(middleware.RequestWrapper(security.NewIpExtractor(true)))

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {

		pkg, err := net.ScrapeInterfacesIO(ctx)
		if err != nil {
			log.Error("scrape error", "error", err)
			api.InternalErrorResponse().Write(w)
			return
		}

		api.NewResponse().
			SetCode(http.StatusOK).
			WrapData(pkg).
			Write(w)
	})

	srv := server.NewServer(m)

	root.NewThread()
	go func() {
		defer root.DoneThread()

		err := srv.Run(ctx, ":3000", "", "")
		if err != nil {
			log.Error("server run error", "error", err)
		}
	}()

	root.WaitThreads(5 * time.Second)
	srv.Close()
}
