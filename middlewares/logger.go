package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"server/config"
	"server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		timeStart := time.Now()
		c.Next()
		timeEnd := time.Now()

		ctx, _ := utils.Ctx(c)
		config := config.GetConfigInstance()

		file, err := os.OpenFile(fmt.Sprintf("%s%s.log", config.LOGGIN_DIR, timeStart.Format("2006-01-02")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			return
		}

		log.SetOutput(file)
		logLine := RequestLogLine{
			Id:            strconv.Itoa(int(ctx.User.Id)),
			TimeIn:        timeStart.Format(time.RFC3339),
			DurationMs:    timeEnd.Sub(timeStart).Milliseconds(),
			StatusCode:    c.Writer.Status(),
			RequestMethod: c.Request.Method,
			RequestURI:    c.Request.RequestURI,
		}

		logLineJson, _ := json.Marshal(logLine)
		log.Printf("%v", string(logLineJson))
	}
}

type RequestLogLine struct {
	Id         string `json:"id"`
	TimeIn     string `json:"time_in"`
	DurationMs int64  `json:"duration_ms"`

	StatusCode    int    `json:"status_code"`
	RequestMethod string `json:"request_method"`
	RequestURI    string `json:"request_url"`
}
