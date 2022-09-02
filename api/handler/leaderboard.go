package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/doorbash/leaderboard-api/api/domain"
	"github.com/doorbash/leaderboard-api/api/repository"
	"github.com/doorbash/leaderboard-api/api/util"
	"github.com/doorbash/leaderboard-api/api/util/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type LeaderboardHandler struct {
	pRepo   domain.PlayerRepository
	ldRepo  domain.LeaderboardDataRepository
	ipCache domain.BannedIpsCache
	ldCache domain.LeaderboardDataCache
	router  *mux.Router
}

func (l *LeaderboardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	lid, ok := mux.Vars(r)["lid"]
	if !ok {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	var offset int
	var count int
	var err error

	offset_str := r.URL.Query().Get("offset")
	if offset_str != "" {
		offset, err = strconv.Atoi(offset_str)
		if err != nil || offset < 0 {
			util.WriteStatus(w, http.StatusBadRequest)
			return
		}
	}

	count_str := r.URL.Query().Get("count")
	if count_str != "" {
		count, err = strconv.Atoi(count_str)
		if err != nil || count < 1 || count > 100 {
			util.WriteStatus(w, http.StatusBadRequest)
			return
		}
	} else {
		count = 100
	}

	if offset == 0 && count == 100 {
		ctx, cancel := util.GetContextWithTimeout(r.Context())
		defer cancel()
		ls, err := l.ldCache.GetByUID(ctx, lid)
		if err != nil && err != redis.Nil {
			log.Println(err)
			util.WriteInternalServerError(w)
			return
		}

		if ls != nil {
			util.WriteJson(w, json.RawMessage(*ls))
			return
		}
	}

	ctx, cancel := util.GetContextWithTimeout(r.Context())
	defer cancel()
	lds, err := l.ldRepo.GetByUID(ctx, lid, offset, count)
	if err != nil {
		log.Println(err)
		util.WriteInternalServerError(w)
		return
	}

	if offset == 0 && count == 100 {
		ctx, cancel := util.GetContextWithTimeout(r.Context())
		defer cancel()
		err = l.ldCache.Set(ctx, lid, lds)
		if err != nil {
			util.WriteInternalServerError(w)
			return
		}
	}

	util.WriteJson(w, lds)
}

func (l *LeaderboardHandler) GetLeaderboardByPID(w http.ResponseWriter, r *http.Request) {
	lid, ok := mux.Vars(r)["lid"]
	if !ok {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	pid, ok := mux.Vars(r)["pid"]
	if !ok {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	var offset int
	var count int
	var err error

	offset_str := r.URL.Query().Get("offset")
	if offset_str != "" {
		offset, err = strconv.Atoi(offset_str)
		if err != nil {
			util.WriteStatus(w, http.StatusBadRequest)
			return
		}
	} else {
		offset = -50
	}

	count_str := r.URL.Query().Get("count")
	if count_str != "" {
		count, err = strconv.Atoi(count_str)
		if err != nil || count < 1 || count > 100 {
			util.WriteStatus(w, http.StatusBadRequest)
			return
		}
	} else {
		count = 100
	}

	ctx, cancel := util.GetContextWithTimeout(r.Context())
	defer cancel()
	pos, err := l.ldRepo.GetPlayerPosition(ctx, lid, pid)
	if err != nil {
		log.Println(err)
		util.WriteInternalServerError(w)
		return
	}

	_offset := 0
	if pos+offset > 0 {
		_offset = pos + offset
	}

	ctx, cancel = util.GetContextWithTimeout(r.Context())
	defer cancel()
	lds, err := l.ldRepo.GetByUID(ctx, lid, _offset, count)
	if err != nil {
		log.Println(err)
		util.WriteInternalServerError(w)
		return
	}

	ret := map[string]interface{}{
		"offset":      _offset,
		"position":    pos,
		"leaderboard": lds,
	}

	util.WriteJson(w, ret)
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
	ldCache domain.LeaderboardDataCache,
) *LeaderboardHandler {
	l := &LeaderboardHandler{
		pRepo:   pRepo,
		ldRepo:  ldRepo,
		ipCache: ipCache,
		ldCache: ldCache,
		router:  r.PathPrefix("/leaderboards").Subrouter(),
	}

	l.router.HandleFunc("/{lid}", l.GetLeaderboard).Methods("GET")
	l.router.HandleFunc("/{lid}/{pid}", l.GetLeaderboardByPID).Methods("GET")

	jsonRouter := l.router.NewRoute().Subrouter()
	jsonRouter.Use(middleware.JsonBodyMiddleware)
	jsonRouter.HandleFunc("/{lid}/new", l.NewLeaderboardData).Methods("POST")

	return l
}
