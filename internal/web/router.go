package web

import (
	"context"
	"net/http"
	"strings"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/services/hostinfo"
	"github.com/eterline/fstmon/internal/web/controller"
	"github.com/eterline/fstmon/internal/web/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRouter(ctx context.Context, cfg config.Configuration) http.Handler {

	hCtrl := controller.NewHostController(
		hostinfo.InitHostInfo(),
		hostinfo.InitCpuLoader(ctx),
	)

	root := chi.NewMux()

	root.Use(
		middleware.RequestWrapper,
		middleware.RequestLogger,
		middleware.NoCacheControl,
		middleware.SecureControl,
	)

	root.NotFound(controller.NotFound)
	root.MethodNotAllowed(controller.BadMethod)

	root.Get("/info", controller.HandleInfo)

	root.With(
		middleware.SourceSubnetsAllow(strings.Join(cfg.AllowedSubnets, " ")),
		middleware.AllowedHosts(strings.Join(cfg.AllowedHosts, " ")),
		middleware.BearerCheck(cfg.AuthToken),
	).Route(
		"/api", func(r chi.Router) {
			r.Get("/net", hCtrl.HandleNetworking)
			r.Get("/sys", hCtrl.HandleSys)
			r.Get("/parts", hCtrl.HandleParts)
			r.Get("/avgload", hCtrl.HandleAvgload)
			r.Get("/temp", hCtrl.HandleTemp)
			r.Get("/cpu", hCtrl.HandleCpu)
		},
	)

	return root
}
