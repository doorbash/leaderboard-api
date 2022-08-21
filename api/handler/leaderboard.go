package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/doorbash/leaderboard-api/api/domain"
	"github.com/doorbash/leaderboard-api/api/repository"
	"github.com/doorbash/leaderboard-api/api/util"
	"github.com/doorbash/leaderboard-api/api/util/middleware"
	"github.com/gorilla/mux"
)

type LeaderboardHandler struct {
	pRepo   domain.PlayerRepository
	ldRepo  domain.LeaderboardDataRepository
	ipCache domain.BannedIpsCache
	router  *mux.Router
}

func (l *LeaderboardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	lid, ok := mux.Vars(r)["lid"]
	if !ok {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	ctx, cancel := util.GetContextWithTimeout(r.Context())
	defer cancel()
	lds, err := l.ldRepo.GetByUID(ctx, lid)
	if err != nil {
		log.Println(err)
		util.WriteInternalServerError(w)
		return
	}
	util.WriteJson(w, lds)
}

func (l *LeaderboardHandler) NewLeaderboardData(w http.ResponseWriter, r *http.Request) {
	lid, ok := mux.Vars(r)["lid"]
	if !ok {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	jsonBody := r.Context().Value("json")

	body, ok := jsonBody.(map[string]interface{})
	if !ok {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	pid, _ := body["pid"].(string)

	name, _ := body["name"].(string)

	value1, _ := body["value1"].(float64)
	v1 := int(value1)

	value2, _ := body["value2"].(float64)
	v2 := int(value2)

	value3, _ := body["value3"].(float64)
	v3 := int(value3)

	var player *domain.Player
	var err error
	var ctx context.Context
	var cancel context.CancelFunc

	if pid != "" {
		ctx, cancel := util.GetContextWithTimeout(r.Context())
		defer cancel()
		player, err = l.pRepo.GetByUID(ctx, pid)

		if err != nil {
			util.WriteInternalServerError(w)
			return
		}

		if player.Banned {
			util.WriteOK(w)
			return
		}

		if name != "" && player.Name != name {
			player.Name = name
			ctx, cancel = util.GetContextWithTimeout(r.Context())
			defer cancel()
			err = l.pRepo.Update(ctx, player)
			if err != nil {
				util.WriteInternalServerError(w)
				return
			}
		}
	} else {
		if name == "" {
			util.WriteStatus(w, http.StatusBadRequest)
			return
		}
		player = &domain.Player{
			Name: name,
		}
		ctx, cancel = util.GetContextWithTimeout(r.Context())
		defer cancel()
		err = l.pRepo.Insert(ctx, player)
		if err != nil {
			util.WriteInternalServerError(w)
			return
		}
	}

	ctx, cancel = util.GetContextWithTimeout(r.Context())
	defer cancel()
	err = l.ldRepo.Insert(ctx, lid, player.UID, v1, v2, v3)

	if err != nil {
		log.Println(err)
		if err == repository.ErrBetterDataAlreadyExists {
			util.WriteStatus(w, http.StatusConflict)
			return
		}
		if err == repository.ErrLimit {
			ip := r.Header.Get("X-Forwarded-For")
			ctx, cancel = util.GetContextWithTimeout(r.Context())
			defer cancel()
			err = l.ipCache.BanThisIP(ctx, ip)
			if err != nil {
				log.Println(err)
				util.WriteInternalServerError(w)
				return
			}

			player.Banned = true
			ctx, cancel = util.GetContextWithTimeout(r.Context())
			defer cancel()
			err = l.pRepo.Update(ctx, player)
			if err != nil {
				util.WriteInternalServerError(w)
				return
			}
			util.WriteOK(w)
			return
		}
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	util.WriteJson(w, *player)
}

func NewLeaderboardHandler(
	r *mux.Router,
	pRepo domain.PlayerRepository,
	ldRepo domain.LeaderboardDataRepository,
	ipCache domain.BannedIpsCache,
) *LeaderboardHandler {
	l := &LeaderboardHandler{
		pRepo:   pRepo,
		ldRepo:  ldRepo,
		ipCache: ipCache,
		router:  r.PathPrefix("/leaderboards").Subrouter(),
	}

	l.router.HandleFunc("/{lid}", l.GetLeaderboard).Methods("GET")

	jsonRouter := l.router.NewRoute().Subrouter()
	jsonRouter.Use(middleware.JsonBodyMiddleware)
	jsonRouter.HandleFunc("/{lid}/new", l.NewLeaderboardData).Methods("POST")

	return l
}
