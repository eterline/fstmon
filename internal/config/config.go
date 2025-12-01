// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package config

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"golang.org/x/sys/unix"
)

func clampSeconds(sec, min, max int) time.Duration {
	if sec < min {
		sec = min
	}
	if sec > max {
		sec = max
	}
	return time.Duration(sec) * time.Second
}

type Monitor struct {
	Cpu       int `arg:"--cpu-loop" help:"Cpu metric update loop seconds"`
	Memory    int `arg:"--memory-loop" help:"Memory update loop seconds"`
	System    int `arg:"--system-loop" help:"System update loop seconds"`
	Thermal   int `arg:"--thermal-loop" help:"Thermal update loop seconds"`
	NetworkIO int `arg:"--network-loop" help:"Network I/O update loop seconds"`
	DiskIO    int `arg:"--partitions-loop" help:"Disk I/O update loop seconds"`
}

func (m Monitor) CpuDuration() time.Duration {
	return clampSeconds(m.Cpu, 5, 120)
}

func (m Monitor) NetworkIoDuration() time.Duration {
	return clampSeconds(m.NetworkIO, 5, 120)
}

func (m Monitor) MemorykDuration() time.Duration {
	return clampSeconds(m.Memory, 5, 120)
}

func (m Monitor) SystemDuration() time.Duration {
	return clampSeconds(m.System, 15, 300)
}

func (m Monitor) ThermalDuration() time.Duration {
	return clampSeconds(m.Thermal, 15, 300)
}

func (m Monitor) DiskIODuration() time.Duration {
	return clampSeconds(m.DiskIO, 10, 300)
}

type (
	Log struct {
		LogLevel      string `arg:"--log-level" help:"Logging level: debug|info|warn|error"`
		JSONlog       bool   `arg:"--log-json,-j" help:"Set logs to JSON format"`
		AccessLogFile string `arg:"--access-log" help:"Set access log file"`
	}

	Server struct {
		Listen     string `arg:"--listen,-l" help:"Server listen address"`
		CrtFileSSL string `arg:"--certfile,-c" help:"Server SSL certificate file"`
		KeyFileSSL string `arg:"--keyfile,-k" help:"Server SSL key file"`
	}

	Secure struct {
		AllowedHosts   []string `arg:"--sni,-h" help:"Server allowed request hosts"`
		AllowedSubnets []string `arg:"--subnets,-s" help:"Server allowed source subnets/IPs"`
		AuthToken      []string `arg:"--token,-t,env:TOKEN" help:"Server auth token string"`
		ParseIpHeader  bool     `arg:"--ip-header" help:"Enable parsing reverse proxy headers"`
	}

	Configuration struct {
		Log
		Server
		Secure
		Monitor
	}
)

func (l Log) AccessLog() (wr io.WriteCloser, err error, wrEnable bool) {
	if strings.ToLower(l.AccessLogFile) == "none" {
		return nil, nil, false
	}

	out := io.WriteCloser(os.Stdout)

	if strings.ToLower(l.AccessLogFile) == "stdout" {
		return out, nil, true
	}

	dir := filepath.Dir(l.AccessLogFile)
	ds, err := os.Lstat(dir)
	if err != nil {
		return out, err, true
	}

	// symlink attack prevent
	if (ds.Mode() & os.ModeSymlink) != 0 {
		return out, err, true
	}

	fd, err := unix.Open(
		l.AccessLogFile,
		unix.O_WRONLY|unix.O_APPEND|unix.O_CREAT|unix.O_NOFOLLOW,
		0644,
	)
	if err != nil {
		return out, err, true
	}

	out = os.NewFile(uintptr(fd), l.AccessLogFile)

	return out, nil, true
}

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
