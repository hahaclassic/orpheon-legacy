package entity

import "github.com/google/uuid"

type Claims struct {
	UserID    uuid.UUID   `json:"user_id"`
	AccessLvl AccessLevel `json:"access_lvl"`
}
