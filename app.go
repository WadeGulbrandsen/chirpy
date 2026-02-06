package main

import (
	"net/http"
)

func appHandler(cfg *apiConfig) http.Handler {
	appDir := http.Dir(cfg.appPath)
	return http.StripPrefix(
		cfg.appPrefix,
		cfg.middlewareMetricsInc(http.FileServer(appDir)),
	)
}
