package config

import (
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

type Configuration struct {
	Debug   bool `arg:"--debug" help:"Allow debug logging level"`
	JSONlog bool `arg:"--log-json" help:"Set logs to JSON format"`

	Listen     string `arg:"--listen" help:"Server listen address"`
	CrtFileSSL string `arg:"--certfile,env:CERT" help:"Server SSL certificate file"`
	KeyFileSSL string `arg:"--keyfile,env:KEY" help:"Server SSL key file"`

	AllowedHosts   []string `arg:"--sni" help:"Server allowed request hosts"`
	AllowedSubnets []string `arg:"--subnets" help:"Server allowed source subnets/IPs"`
	AuthToken      string   `arg:"--token" help:"Server auth token string"`
	ParseIpHeader  bool     `arg:"--ip-header" help:"Enable parsing reverse proxy headers"`
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
