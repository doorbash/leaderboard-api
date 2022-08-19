package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/doorbash/leaderboard-api/api/domain"
	"github.com/doorbash/leaderboard-api/api/util"
	"github.com/gorilla/mux"
)

type LevelHandler struct {
	lRepo  domain.LevelRepository
	router *mux.Router
}

func (l *LevelHandler) GetAllLevels(w http.ResponseWriter, r *http.Request) {
	gid, ok := mux.Vars(r)["gid"]
	if !ok {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	ctx, cancel := util.GetContextWithTimeout(r.Context())
	defer cancel()
	levels, err := l.lRepo.GetAllByGameID(ctx, gid)
	if err != nil {
		log.Println(err)
		util.WriteInternalServerError(w)
		return
	}
	util.WriteJson(w, levels)
}

func (l *LevelHandler) GetLevel(w http.ResponseWriter, r *http.Request) {
	gid, ok := mux.Vars(r)["gid"]
	if !ok {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	n, ok := mux.Vars(r)["number"]
	if !ok {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	number, err := strconv.Atoi(n)
	if err != nil {
		util.WriteStatus(w, http.StatusBadRequest)
		return
	}

	ctx, cancel := util.GetContextWithTimeout(r.Context())
	defer cancel()
	level, err := l.lRepo.GetByGameIDAndNumber(ctx, gid, number)
	if err != nil {
		log.Println(err)
		util.WriteInternalServerError(w)
		return
	}
	util.WriteJson(w, level)
}

func NewLevelHandler(r *mux.Router, lRepo domain.LevelRepository) *LevelHandler {
	l := &LevelHandler{
		lRepo:  lRepo,
		router: r.NewRoute().Subrouter(),
	}

	l.router.HandleFunc("/{gid}/levels", l.GetAllLevels).Methods("GET")
	l.router.HandleFunc("/{gid}/{number}", l.GetLevel).Methods("GET")

	return l
}
