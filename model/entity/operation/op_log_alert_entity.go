package operation

import (
	"time"

	"github.com/google/uuid"
)

type LogAlert struct {
	Id         uuid.UUID
	TruckId    uuid.UUID
	AlertType  string
	AiAnalysis string
	IsResolved bool
	CreatedAt  time.Time
}

func (b *LogAlert) TableName() string {
	return "operations.log_alerts"
}
