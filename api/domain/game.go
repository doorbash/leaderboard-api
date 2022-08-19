package domain

import (
	"context"
)

type Game struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	VersionName string `json:"version_name"`
}

type GameRepository interface {
	GetAll(ctx context.Context) ([]Game, error)
	GetByID(ctx context.Context, id string) (*Game, error)
}
