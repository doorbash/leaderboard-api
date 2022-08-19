package handler

import (
	"log"
	"net/http"

	"github.com/doorbash/leaderboard-api/api/domain"
	"github.com/doorbash/leaderboard-api/api/util"
	"github.com/gorilla/mux"
)

type GameHandler struct {
	gRepo  domain.GameRepository
	router *mux.Router
}

func (g *GameHandler) GetAllGamesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := util.GetContextWithTimeout(r.Context())
	defer cancel()
	games, err := g.gRepo.GetAll(ctx)
	if err != nil {
		log.Println(err)
		util.WriteInternalServerError(w)
		return
	}
	util.WriteJson(w, games)
}

func (g *GameHandler) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	gid, ok := mux.Vars(r)["id"]
	if !ok {
		util.WriteInternalServerError(w)
		return
	}
	ctx, cancel := util.GetContextWithTimeout(r.Context())
	defer cancel()
	game, err := g.gRepo.GetByID(ctx, gid)
	if err != nil {
		util.WriteInternalServerError(w)
		return
	}
	util.WriteJson(w, game)
}

func NewGameHandler(r *mux.Router, gRepo domain.GameRepository) *GameHandler {
	g := &GameHandler{
		gRepo:  gRepo,
		router: r.NewRoute().Subrouter(),
	}

	g.router.HandleFunc("/games", g.GetAllGamesHandler).Methods("GET")
	g.router.HandleFunc("/{id}/", g.GetGameHandler).Methods("GET")

	return g
}
