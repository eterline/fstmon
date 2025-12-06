// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package main

import (
	"fmt"
	"os"

	"github.com/eterline/fstmon/internal/infra/xrayapi"
	"github.com/eterline/fstmon/internal/log"
)

func main() {
	lg := log.NewLogger("debug", false)
	// app.Execute(root, Flags, cfg)

	lg.Info("xray api init")

	api, err := xrayapi.New("10.192.0.52:3000")
	if err != nil {
		lg.Error("xray api conn error", "error", err.Error())
		os.Exit(0)
	}
	defer api.Close()

	fmt.Println(api.GetTraffic(false))

}
