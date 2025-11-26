// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package app

import (
	"os"
	"time"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/log"
	"github.com/eterline/fstmon/internal/web"
	"github.com/eterline/fstmon/internal/web/server"
	"github.com/eterline/fstmon/pkg/toolkit"
)

type InitFlags struct {
	CommitHash string
	Version    string
}

func Execute(root *toolkit.AppStarter, flags InitFlags, cfg config.Configuration) {
	log := log.MustLoggerFromContext(root.Context)

	log.Info("app started", "commit", flags.CommitHash, "version", flags.Version)
	defer func() {
		log.Info("app closed", "working_time", root.WorkTime())
	}()

	routes := web.RegisterRouter(root.Context, cfg)
	srv := server.NewServer(routes)
	defer srv.Close()

	root.NewThread()
	go func() {
		defer root.DoneThread()
		err := srv.Run(root.Context, cfg.Listen, cfg.KeyFileSSL, cfg.CrtFileSSL)
		if err != nil {
			log.Error("server exited with error", "error", err)
			root.StopApp()
		}
	}()

	if err := root.WaitThreads(5 * time.Second); err != nil {
		log.Warn("force exit")
		os.Exit(1)
	}
}
