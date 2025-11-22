package entity

type User struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	TeamName string `db:"team_name"`
	IsActive bool   `db:"is_active"`
}

type UserItem struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
