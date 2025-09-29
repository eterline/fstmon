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

var (
	cfg = config.Configuration{
		Log: config.Log{
			Debug:   false,
			JSONlog: true,
		},
		Server: config.Server{
			Listen:     ":3000",
			CrtFileSSL: "",
			KeyFileSSL: "",
		},
		Secure: config.Secure{
			AllowedSubnets: []string{},
			AllowedHosts:   []string{},
			AuthToken:      "",
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

	log.InitLogger(cfg.Debug, cfg.JSONlog)
	app.Execute(root, cfg)
}
