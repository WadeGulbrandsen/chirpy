package main

import "net/http"

const (
	app_path  = "."
	appDir    = http.Dir(app_path)
	appPrefix = "/app/"
)

func appHandler(apiCfg *apiConfig) http.Handler {
	return http.StripPrefix(
		appPrefix,
		apiCfg.middlewareMetricsInc(http.FileServer(appDir)),
	)
}
