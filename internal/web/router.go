package web

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/internal/services/hostfetchers"
	"github.com/eterline/fstmon/internal/services/monitors"
	"github.com/eterline/fstmon/internal/web/controller"
	"github.com/eterline/fstmon/internal/web/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRouter(ctx context.Context, cfg config.Configuration) http.Handler {

	root := chi.NewMux()

	root.Use(
		middleware.RequestWrapper,
		middleware.RequestLogger,
		middleware.NoCacheControl,
		middleware.SecureControl,
	)

	root.NotFound(controller.NotFound)
	root.MethodNotAllowed(controller.BadMethod)

	root.With(
		middleware.SourceSubnetsAllow(strings.Join(cfg.AllowedSubnets, " ")),
		middleware.AllowedHosts(strings.Join(cfg.AllowedHosts, " ")),
		middleware.BearerCheck(cfg.AuthToken),
	).Route(
		"/monitoring", func(r chi.Router) {

			cpu5 := monitors.InitCpuLoadMon(ctx, 5*time.Second)

			hc := controller.NewHostController(
				hostfetchers.InitSystemMon(ctx, cpu5),
				hostfetchers.InitAverageLoadMon(ctx),
				hostfetchers.InitPartUseMon(ctx),
				hostfetchers.InitNetworkMon(ctx),
			)

			r.Get("/system", hc.HandleSystem)
			r.Get("/avgload", hc.HandleAvgload)
			r.Get("/parts", hc.HandleParts)
			r.Get("/net", hc.HandleNetworking)
		},
	)

	return root
}
