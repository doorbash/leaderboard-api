package handler

import (
	"log"
	"net/http"

	"github.com/doorbash/leaderboard-api/api/domain"
	"github.com/doorbash/leaderboard-api/api/util"
	"github.com/gorilla/mux"
)

type LeaderboardHandler struct {
	pRepo  domain.PlayerRepository
	ldRepo domain.LeaderboardDataRepository
	router *mux.Router
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

func NewLeaderboardHandler(
	r *mux.Router,
	pRepo domain.PlayerRepository,
	ldRepo domain.LeaderboardDataRepository,
) *LeaderboardHandler {
	l := &LeaderboardHandler{
		pRepo:  pRepo,
		ldRepo: ldRepo,
		router: r.NewRoute().Subrouter(),
	}

	l.router.HandleFunc("/leaderboards/{lid}", l.GetLeaderboard).Methods("GET")

	return l
}
