package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/doorbash/leaderboard-api/api/domain"
	"github.com/doorbash/leaderboard-api/api/util"
)

type BannedIPHandler struct {
	IPsCache domain.BannedIpsCache
}

func (b BannedIPHandler) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")

		ctx, cancel := util.GetContextWithTimeout(context.Background())
		defer cancel()
		banned, err := b.IPsCache.IsThisIPBanned(ctx, ip)

		if err != nil {
			log.Println(err)
			util.WriteInternalServerError(w)
			return
		}

		if banned {
			util.WriteInternalServerError(w)
			return
		}

		h.ServeHTTP(w, r)
	})
}
