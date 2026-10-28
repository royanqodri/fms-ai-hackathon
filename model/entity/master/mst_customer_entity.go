package master

import (
	"time"

	"github.com/google/uuid"
)

type MstCustomer struct {
	Id         uuid.UUID
	CustomerNo string
	Name       string
	Type       string
	PhoneNo    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (b *MstCustomer) TableName() string {
	return "master.mst_customers"
}
