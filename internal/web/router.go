package web

import (
	"net/http"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/services/hostinfo"
	"github.com/eterline/fstmon/internal/web/controller"
	"github.com/eterline/fstmon/internal/web/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRouter(cfg config.Configuration) http.Handler {

	hCtrl := controller.NewHostController(
		hostinfo.InitHostInfo(cfg.Interface),
	)

	root := chi.NewMux()

	root.With(
		middleware.RequestWrapper,
		middleware.RequestLogger,
		middleware.NoCacheControl,
		middleware.SecureControl,
		middleware.SourceSubnetsAllow(cfg.AllowedSubnets),
	).Route("/api", func(r chi.Router) {

		r.Get("/net", hCtrl.HandleNetworking)
		r.Get("/sys", hCtrl.HandleSys)
		r.Get("/proc", hCtrl.HandleSys)

	})

	return root
}
