package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"system.buon18.com/m/config"
)

type WrapperResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Size       int
}

func NewWrapperResponseWriter(w http.ResponseWriter) *WrapperResponseWriter {
	return &WrapperResponseWriter{w, http.StatusOK, 0}
}

func (w *WrapperResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *WrapperResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.Size += size
	return size, err
}

func LoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		ww := NewWrapperResponseWriter(w)
		next.ServeHTTP(ww, r)
		timeEnd := time.Now()

		config := config.GetConfigInstance()

		if config.LOGGING_DIR == "" {
			return
		}

		path := filepath.Join(config.LOGGING_DIR, fmt.Sprintf("%s.log", timeStart.Format("2006-01-02")))

		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.MkdirAll(config.LOGGING_DIR, os.ModePerm)
			if err != nil {
				slog.Error(fmt.Sprintf("Error creating directory: %v", err))
				return
			}
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		defer func() {
			if err := file.Close(); err != nil {
				slog.Error(fmt.Sprintf("Error closing file: %v", err))
			}
		}()
		if err != nil {
			slog.Error(fmt.Sprintf("Error opening file: %v", err))
			return
		}

		log.SetOutput(file)
		logLine := requestLogLine{
			ID:            strconv.Itoa(0),
			TimeIn:        timeStart.Format(time.RFC3339),
			DurationMs:    timeEnd.Sub(timeStart).Milliseconds(),
			StatusCode:    ww.StatusCode,
			RequestMethod: r.Method,
			RequestURI:    r.RequestURI,
		}

		logLineJson, _ := json.Marshal(logLine)
		log.Printf("%v", string(logLineJson))
	})
}

type requestLogLine struct {
	ID         string `json:"id"`
	TimeIn     string `json:"time_in"`
	DurationMs int64  `json:"duration_ms"`

	StatusCode    int    `json:"status_code"`
	RequestMethod string `json:"request_method"`
	RequestURI    string `json:"request_url"`
}
