// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package main

import (
	"fmt"
	"os"

	"github.com/eterline/fstmon/internal/app"
	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/infra/log"
	"github.com/eterline/fstmon/internal/infra/security"
	"github.com/eterline/fstmon/pkg/toolkit"
)

func init() {
	toolkit.RegisterCommand("newtoken", func(args ...string) error {
		ti, err := security.NewIssuedTokenAuthProvide(true)
		if err != nil {
			return err
		}
		token, err := ti.Issue()
		if err != nil {
			return err
		}
		fmt.Printf("token: %s\n", token)
		return nil
	})
}

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
			LogLevel:      "info",
			JSONlog:       false,
			AccessLogFile: "stdout",
		},
		Server: config.Server{
			Listen:     ":3000",
			CrtFileSSL: "",
			KeyFileSSL: "",
		},
		Secure: config.Secure{
			AllowedSubnets: []string{},
			AllowedHosts:   []string{},
			AuthToken:      false,
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
			err, ok := app.RunAdditional()
			if err != nil {
				return err
			}

			if ok {
				os.Exit(0)
			}

			err = config.ParseArgs(&cfg)
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
