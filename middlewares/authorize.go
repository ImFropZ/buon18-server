package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"system.buon18.com/m/utils"
)

func Authorize(next http.Handler, allowPermissions []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)

		if userCtx == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusUnauthorized, "missing 'Authorization' header or 'Authorization' header's value doesn't start with 'Bearer '", "Unauthorized", nil))
			return
		}

		// -- Add FULL_ACCESS permission
		allowPermissions = append(allowPermissions, "FULL_ACCESS")

		allow := false
		for _, permission := range allowPermissions {
			// -- Check permission
			for _, ctxPermission := range *userCtx.Permissions {
				if strings.EqualFold(ctxPermission.Name, permission) {
					allow = true
					break
				}
			}
		}

		if !allow {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusForbidden, "", "Forbidden", nil))
			return
		}

		next.ServeHTTP(w, r)
	})
}
