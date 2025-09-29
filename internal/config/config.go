// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package config

import (
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

type (
	Log struct {
		Debug   bool `arg:"--debug,-d" help:"Allow debug logging level"`
		JSONlog bool `arg:"--log-json,-j" help:"Set logs to JSON format"`
	}

	Server struct {
		Listen     string `arg:"--listen,-l" help:"Server listen address"`
		CrtFileSSL string `arg:"--certfile,-c,env:CERT" help:"Server SSL certificate file"`
		KeyFileSSL string `arg:"--keyfile,-k,env:KEY" help:"Server SSL key file"`
	}

	Secure struct {
		AllowedHosts   []string `arg:"--sni,-h" help:"Server allowed request hosts"`
		AllowedSubnets []string `arg:"--subnets,-s" help:"Server allowed source subnets/IPs"`
		AuthToken      string   `arg:"--token,-t,env:TOKEN" help:"Server auth token string"`
		ParseIpHeader  bool     `arg:"--ip-header" help:"Enable parsing reverse proxy headers"`
	}

	Configuration struct {
		Log
		Server
		Secure
	}
)

var (
	parserConfig = arg.Config{
		Program:           selfExec(),
		IgnoreEnv:         false,
		IgnoreDefault:     false,
		StrictSubcommands: true,
	}
)

func ParseArgs(c *Configuration) error {
	p, err := arg.NewParser(parserConfig, c)
	if err != nil {
		return err
	}

	err = p.Parse(os.Args[1:])
	if err == arg.ErrHelp {
		p.WriteHelp(os.Stdout)
		os.Exit(1)
	}
	return err
}

func selfExec() string {
	exePath, err := os.Executable()
	if err != nil {
		return "monita"
	}

	return filepath.Base(exePath)
}
