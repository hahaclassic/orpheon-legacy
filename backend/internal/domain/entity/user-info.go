package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserInfo struct {
	ID               uuid.UUID   `json:"id"`
	Name             string      `json:"name"`
	RegistrationDate time.Time   `json:"registration_date"`
	BirthDate        time.Time   `json:"birth_date"`
	AccessLvl        AccessLevel `json:"access_lvl"`
}
