package middlewares

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/utils"

	"github.com/nullism/bqb"
)

func Authenticate(next http.Handler, DB *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// -- Get token
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusUnauthorized, "missing 'Authorization' header or 'Authorization' header's value doesn't start with 'Bearer '", "Unauthorized", nil))
			return
		}

		claims, err := utils.ValidateWebToken(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusUnauthorized, "invalid token", "Unauthorized", nil))
			return
		}

		// -- Prepare sql query
		query, params, err := bqb.New(`
		SELECT 
			"setting.user".id, 
			"setting.user".name, 
			"setting.user".email, 
			"setting.user".typ, 
			COALESCE("setting.role".id, 0), 
			COALESCE("setting.role".name, ''), 
			COALESCE("setting.role".description, ''), 
			COALESCE("setting.permission".id, 0), 
			COALESCE("setting.permission".name, '')
		FROM 
			"setting.user"
		LEFT JOIN "setting.role" ON "setting.user".setting_role_id = "setting.role".id
		LEFT JOIN "setting.role_permission" ON "setting.role".id = "setting.role_permission".setting_role_id 
		LEFT JOIN "setting.permission" ON "setting.role_permission".setting_permission_id = "setting.permission".id
		WHERE "setting.user".email = ?
		ORDER BY "setting.user".email, "setting.role".id, "setting.permission".id`, claims.Email).ToPgsql()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusUnauthorized, "internal server error", "Unauthorized", nil))
			return
		}

		// -- Validate user
		rows, err := DB.Query(query, params...)
		if err != nil {
			slog.Error(fmt.Sprintf("Error querying user: %v", err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusInternalServerError, "internal server error", "Internal Server Error", nil))
			return
		}

		var user setting.SettingUser
		var role setting.SettingRole
		permissions := make([]setting.SettingPermission, 0)
		for rows.Next() {
			var permission setting.SettingPermission
			err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Typ, &role.Id, &role.Name, &role.Description, &permission.Id, &permission.Name)
			if err != nil {
				slog.Error(fmt.Sprintf("Error scanning user: %v", err))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusInternalServerError, "internal server error", "Internal Server Error", nil))
				return
			}

			permissions = append(permissions, permission)
		}

		if user.Id == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusUnauthorized, "user not found", "Unauthorized", nil))
			return
		}

		userCtx := utils.CtxValue{
			User:        &user,
			Role:        &role,
			Permissions: &permissions,
		}

		r = r.WithContext(context.WithValue(r.Context(), utils.CtxKey{}, &userCtx))
		next.ServeHTTP(w, r)
	})
}
