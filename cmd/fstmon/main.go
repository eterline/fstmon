package main

import (
	"github.com/eterline/fstmon/internal/app"
	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/log"
	"github.com/eterline/fstmon/pkg/toolkit"
)

var (
	cfg = config.Configuration{
		Debug:          false,
		Listen:         ":8100",
		CrtFileSSL:     "",
		KeyFileSSL:     "",
		AllowedSubnets: "",
	}
)

func main() {

	log.InitLogger(false)

	root := toolkit.InitAppStart(
		func() error {
			err := config.ParseArgs(&cfg)
			if err != nil {
				return err
			}
			return nil
		},
	)

	app.Execute(root, cfg)
}
