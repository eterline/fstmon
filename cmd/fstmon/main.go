// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package main

import (
	"github.com/eterline/fstmon/internal/app"
	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/log"
	"github.com/eterline/fstmon/pkg/toolkit"
)

// -ladflags variables
var (
	CommitHash = "dev"
	Version    = "dev"
)

var (
	Flags = app.InitFlags{
		CommitHash: CommitHash,
		Version:    Version,
	}

	cfg = config.Configuration{
		Log: config.Log{
			LogLevel:        "info",
			JSONlog:         false,
			AccessLogFile:   "",
			AccessLogEnable: false,
		},
		Server: config.Server{
			Listen:     ":3000",
			CrtFileSSL: "",
			KeyFileSSL: "",
		},
		Secure: config.Secure{
			AllowedSubnets: []string{},
			AllowedHosts:   []string{},
			AuthToken:      []string{},
		},
		Monitor: config.Monitor{
			Cpu:       10,
			Memory:    10,
			System:    20,
			Thermal:   20,
			NetworkIO: 10,
			DiskIO:    10,
		},
	}
)

func main() {
	root := toolkit.InitAppStart(
		func() error {
			err := config.ParseArgs(&cfg)
			if err != nil {
				return err
			}
			return nil
		},
	)

	logger := log.NewLogger(cfg.LogLevel, cfg.JSONlog)
	root.Context = log.WrapLoggerToContext(root.Context, logger)

	app.Execute(root, Flags, cfg)
}
