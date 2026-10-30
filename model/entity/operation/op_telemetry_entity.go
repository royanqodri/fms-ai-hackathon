package operation

import (
	"time"

	"github.com/google/uuid"
)

type TrxTelemetry struct {
	Id          uuid.UUID
	TruckId     uuid.UUID
	Latitude    float64
	Longitude   float64
	Temperature float64
	WeightLoad  float64
	Vibration   float64
	FuelLevel   float64
	OilPressure float64
	EngineRpm   float64
	RecordedAt  time.Time
}

func (b *TrxTelemetry) TableName() string {
	return "operations.trx_telemetries"
}
