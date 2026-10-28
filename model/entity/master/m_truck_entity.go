package master

import (
	"time"

	"github.com/google/uuid"
)

type MstTruck struct {
	Id          uuid.UUID
	CustomerId  uuid.UUID
	PlateNumber string
	CargoType   string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (b *MstTruck) TableName() string {
	return "master.mst_trucks"
}
