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
		hostinfo.InitHostInfo(),
	)

	root := chi.NewMux()
	root.Use(
		middleware.RequestWrapper,
		middleware.RequestLogger,
		middleware.NoCacheControl,
		middleware.SecureControl,
		middleware.SourceSubnetsAllow(cfg.AllowedSubnets),
		middleware.BearerCheck(cfg.AuthToken),
	)

	root.Route("/api", func(r chi.Router) {
		r.Get("/net", hCtrl.HandleNetworking)
		r.Get("/sys", hCtrl.HandleSys)
		r.Get("/parts", hCtrl.HandleParts)
		r.Get("/avgload", hCtrl.HandleAvgload)
		r.Get("/temp", hCtrl.HandleTemp)
	})

	return root
}
