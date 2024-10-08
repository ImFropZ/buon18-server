package middlewares

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"system.buon18.com/m/utils"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Warn("Recovered from panic: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusInternalServerError, "", "Internal Server Error", nil))
			}
		}()
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
