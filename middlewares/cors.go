package middlewares

import (
	"net/http"

	"system.buon18.com/m/config"
)

func CORSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := config.GetConfigInstance()

		w.Header().Add("Access-Control-Allow-Origin", config.ACCESS_CONTROL_ALLOW_ORIGIN)
		w.Header().Add("Access-Control-Allow-Credentials", config.ACCESS_CONTROL_ALLOW_CREDENTIALS)
		w.Header().Add("Access-Control-Allow-Headers", config.ACCESS_CONTROL_ALLOW_HEADERS)
		w.Header().Add("Access-Control-Allow-Methods", config.ACCESS_CONTROL_ALLOW_METHODS)
		w.Header().Add("Access-Control-Expose-Headers", config.ACCESS_CONTROL_EXPOSE_HEADERS)
		w.Header().Add("Access-Control-Max-Age", config.ACCESS_CONTROL_MAX_AGE)

		if r.Method == "OPTIONS" {
			w.Header().Add("Allow", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}
