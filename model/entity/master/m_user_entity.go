package master

import (
	"time"

	"github.com/google/uuid"
)

type MstUser struct {
	Id         uuid.UUID
	CustomerId uuid.UUID
	TelegramId string
	Name       string
	Role       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (b *MstUser) TableName() string {
	return "master.mst_users"
}
