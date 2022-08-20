package middleware

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/doorbash/leaderboard-api/api/util"
)

func JsonBodyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			util.WriteInternalServerError(w)
			return
		}
		log.Println(string(data))
		var jsonBody interface{}
		err = json.Unmarshal(data, &jsonBody)
		if err != nil {
			log.Println(err)
			util.WriteStatus(w, http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "json", jsonBody)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
