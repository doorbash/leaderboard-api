package domain

type Leaderboard struct {
	ID          int    `json:"id"`
	UID         string `json:"uid"`
	GID         string `json:"gid"`
	Name        string `json:"name"`
	Value1Order int8   `json:"v1_order"`
	Value2Order int8   `json:"v2_order"`
	Value3Order int8   `json:"v3_order"`
}
