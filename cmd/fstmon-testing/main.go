package main

import (
	"github.com/eterline/fstmon/internal/app"
	"github.com/eterline/fstmon/internal/config"
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
			LogLevel: "info",
			JSONlog:  false,
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
		Monitor: config.Monitor{
			Cpu:        5,
			Avgload:    10,
			System:     30,
			Network:    5,
			Partitions: 30,
		},
	}
)

func main() {

}
