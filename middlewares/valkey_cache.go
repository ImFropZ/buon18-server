package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"system.buon18.com/m/config"
	"system.buon18.com/m/database"
	"system.buon18.com/m/utils"

	"github.com/gin-gonic/gin"
)

func ValkeyCache[T interface{}](connection *database.Connection, fieldName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		if valkeyClient := connection.Valkey; valkeyClient == nil {
			c.Next()
			return
		} else {
			if c.Request.URL.Path == c.Request.RequestURI {
				if resourceStr, err := (*valkeyClient).Do(ctx, (*valkeyClient).B().Get().Key(c.Request.RequestURI).Build()).ToString(); err == nil {
					if totalStr, err := (*valkeyClient).Do(ctx, (*valkeyClient).B().Get().Key(fmt.Sprintf("total_%s", c.Request.RequestURI)).Build()).ToString(); err == nil {
						var jsonResponse T
						json.Unmarshal([]byte(resourceStr), &jsonResponse)
						c.Header("X-Cache", "true")
						c.Header("X-Total-Count", totalStr)
						c.JSON(200, utils.NewResponse(200, "", gin.H{
							fieldName: jsonResponse,
						}))
						c.Abort()
						return
					}
				}
			}

			c.Next()

			go func() {
				if c.Request.URL.Path == c.Request.RequestURI {
					config := config.GetConfigInstance()
					if value, ok := c.Get("response"); ok {
						result := value.([]byte)
						resultStr := string(result)
						err := (*valkeyClient).Do(
							ctx,
							(*valkeyClient).
								B().
								Set().
								Key(c.Request.RequestURI).
								Value(resultStr).
								ExatTimestamp(time.Now().Add(time.Duration(config.CACHE_DURATION_SEC)*time.Second).Unix()).
								Build(),
						).Error()
						if err != nil {
							log.Printf("ValkeyCache: %v\n", err)
						}
					}
					if value, ok := c.Get("total"); ok {
						result := utils.IntToStr(value.(int))
						err := (*valkeyClient).Do(
							ctx,
							(*valkeyClient).
								B().
								Set().
								Key(fmt.Sprintf("total_%s", c.Request.RequestURI)).
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
	}
}
