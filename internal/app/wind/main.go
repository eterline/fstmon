package wind

import (
	"github.com/eterline/fstmon/internal/config"
	"github.com/eterline/fstmon/pkg/toolkit"
)

type InitFlags struct {
	CommitHash string
	Version    string
}

func Execute(root *toolkit.AppStarter, flags InitFlags, cfg config.Configuration) {

}
