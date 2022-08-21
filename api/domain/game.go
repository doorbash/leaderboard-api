package domain

import (
	"context"
	"encoding/json"
)

type Game struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
}

func (r *Game) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})
	ret["id"] = r.ID
	ret["name"] = r.Name
	ret["data"] = json.RawMessage(r.Data)
	b, err := json.Marshal(ret)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type GameRepository interface {
	GetByID(ctx context.Context, id string) (*Game, error)
}

type GameCache interface {
	GetByID(ctx context.Context, id string) (*string, error)
	Set(ctx context.Context, game Game) error
}
