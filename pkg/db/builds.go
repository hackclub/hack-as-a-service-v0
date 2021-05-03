package db

import (
	"time"

	"github.com/google/uuid"
)

type Build struct {
	ID        uint      `gorm:"primary_key"`
	ExecID    uuid.UUID `gorm:"index,unique"`
	AppID     uint      `gorm:"index"`
	StartedAt time.Time
	EndedAt   time.Time
	Running   bool
	Events    []string
	Status    int
}
