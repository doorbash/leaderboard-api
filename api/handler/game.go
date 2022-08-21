package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/doorbash/leaderboard-api/api/domain"
	"github.com/doorbash/leaderboard-api/api/util"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type GameHandler struct {
	gRepo  domain.GameRepository
	gCache domain.GameCache
	router *mux.Router
}

func (g *GameHandler) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	gid, ok := mux.Vars(r)["id"]
	if !ok {
		util.WriteInternalServerError(w)
		return
	}

	ctx, cancel := util.GetContextWithTimeout(r.Context())
	defer cancel()
	gs, err := g.gCache.GetByID(ctx, gid)
	if err != nil && err != redis.Nil {
		log.Println(err)
		util.WriteInternalServerError(w)
		return
	}

	if gs != nil {
		util.WriteJson(w, json.RawMessage(*gs))
		return
	}

	ctx, cancel = util.GetContextWithTimeout(r.Context())
	defer cancel()
	game, err := g.gRepo.GetByID(ctx, gid)
	if err != nil {
		log.Println(err)
		util.WriteStatus(w, http.StatusNotFound)
		return
	}

	ctx, cancel = util.GetContextWithTimeout(r.Context())
	defer cancel()
	err = g.gCache.Set(ctx, *game)
	if err != nil {
		log.Println(err)
		util.WriteInternalServerError(w)
		return
	}

	util.WriteJson(w, game)
}

func NewGameHandler(
	r *mux.Router,
	gRepo domain.GameRepository,
	gCache domain.GameCache,
) *GameHandler {
	g := &GameHandler{
		gRepo:  gRepo,
		gCache: gCache,
		router: r.NewRoute().Subrouter(),
	}

	g.router.HandleFunc("/{id}", g.GetGameHandler).Methods("GET")

	return g
}
