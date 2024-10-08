package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"system.buon18.com/m/config"
	"system.buon18.com/m/database"
	"system.buon18.com/m/utils"
)

func ValkeyCache[T interface{}](next http.Handler, connection *database.Connection, fieldName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if valkeyClient := connection.Valkey; valkeyClient == nil {
			next.ServeHTTP(w, r)
			return
		} else {
			if r.URL.Path == r.RequestURI {
				if resourceStr, err := (*valkeyClient).Do(ctx, (*valkeyClient).B().Get().Key(r.RequestURI).Build()).ToString(); err == nil {
					if totalStr, err := (*valkeyClient).Do(ctx, (*valkeyClient).B().Get().Key(fmt.Sprintf("total_%s", r.RequestURI)).Build()).ToString(); err == nil {
						var jsonResponse T
						json.Unmarshal([]byte(resourceStr), &jsonResponse)
						w.Header().Add("X-Cache", "true")
						w.Header().Add("X-Total-Count", totalStr)
						json.NewEncoder(w).Encode(utils.NewResponse(200, "", map[string]interface{}{
							fieldName: jsonResponse,
						}))
						return
					}
				}
			}

			next.ServeHTTP(w, r)

			go func() {
				if r.URL.Path == r.RequestURI {
					config := config.GetConfigInstance()
					if value := r.Context().Value(""); value != nil {
						result := value.([]byte)
						resultStr := string(result)
						err := (*valkeyClient).Do(
							ctx,
							(*valkeyClient).
								B().
								Set().
								Key(r.RequestURI).
								Value(resultStr).
								ExatTimestamp(time.Now().Add(time.Duration(config.CACHE_DURATION_SEC)*time.Second).Unix()).
								Build(),
						).Error()
						if err != nil {
							log.Printf("ValkeyCache: %v\n", err)
						}
					}
					if value := r.Context().Value(""); value != nil {
						result := utils.IntToStr(value.(int))
						err := (*valkeyClient).Do(
							ctx,
							(*valkeyClient).
								B().
								Set().
								Key(fmt.Sprintf("total_%s", r.RequestURI)).
								Value(result).
								ExatTimestamp(time.Now().Add(time.Duration(config.CACHE_DURATION_SEC)*time.Second).Unix()).
								Build(),
						).Error()
						if err != nil {
							log.Printf("ValkeyCache: %v\n", err)
						}
					}
				}
			}()
		}
	},
	)
}
