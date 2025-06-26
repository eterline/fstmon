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

	srv := server.NewServer(web.RegisterRouter(cfg))
	defer srv.Close()

	slog.Info("staring server", "address", cfg.Listen)

	go func() {
		err := srv.Run(cfg.Listen, cfg.KeyFileSSL, cfg.CrtFileSSL)
		if err != nil {
			slog.Error(err.Error())
			root.StopApp()
		}
	}()

	root.Wait()
}
