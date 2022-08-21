package middleware

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func LoggerMiddleware(h http.Handler) http.Handler {
	return handlers.CustomLoggingHandler(os.Stdout, h, func(writer io.Writer, params handlers.LogFormatterParams) {
		log.Println("HttpLogger:", params.Request.Method, params.Request.URL.Path, params.StatusCode, params.Request.Header.Get("X-Forwarded-For"))
	})
}
