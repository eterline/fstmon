// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package app

import (
	"log/slog"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/web"
	"github.com/eterline/fstmon/internal/web/server"
	"github.com/eterline/fstmon/pkg/toolkit"
)

func Execute(root *toolkit.AppStarter, cfg config.Configuration) {

	slog.Info("app started")
	defer slog.Info("app closed")

	routes := web.RegisterRouter(root.Context, cfg)
	srv := server.NewServer(routes)
	defer srv.Close()

	slog.Info("staring server", "address", cfg.Listen)

	go func() {
		err := srv.Run(cfg.Listen, cfg.KeyFileSSL, cfg.CrtFileSSL)
		if err != nil {
			slog.Error(err.Error())
			root.StopApp()
		}
		slog.Info("server closed")
	}()

	root.Wait()
}
