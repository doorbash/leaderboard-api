package domain

import (
	"context"
	"encoding/json"
)

type Level struct {
	GameID        string `json:"gid"`
	Number        int    `json:"number"`
	Data          string `json:"data"`
	MinValidTime  int
	MaxValidScore int
}

func (l *Level) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})
	ret["number"] = l.Number
	ret["data"] = json.RawMessage(l.Data)
	b, err := json.Marshal(ret)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type LevelRepository interface {
	GetAllByGameID(ctx context.Context, gid string) ([]Level, error)
	GetByGameIDAndNumber(ctx context.Context, gid string, number int) (*Level, error)
}
