package operation

import (
	"github.com/gin-gonic/gin"
	dbConn "github.com/royanqodri/fms-ai-hackathon/database"
	"gorm.io/gorm"

	entityOp "github.com/royanqodri/fms-ai-hackathon/model/entity/operation"
	requestOp "github.com/royanqodri/fms-ai-hackathon/model/request/operation"
	responseOp "github.com/royanqodri/fms-ai-hackathon/model/response/operation"
)

type TrxTelemetryRepository interface {
	GetByParams(ctx *gin.Context, tx *gorm.DB, request requestOp.TrxTelemetryGetRequest) (respData []responseOp.TrxTelemetryGetResponse, err error)
	Save(ctx *gin.Context, tx *gorm.DB, req []entityOp.TrxTelemetry) error
}

type TrxTelemetryRepositoryImpl struct{}

func NewTrxTelemetryRepository() TrxTelemetryRepository {
	return &TrxTelemetryRepositoryImpl{}
}

func (repo TrxTelemetryRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return dbConn.DBConn
	}
	return tx
}

func (repo TrxTelemetryRepositoryImpl) GetByParams(ctx *gin.Context, tx *gorm.DB, request requestOp.TrxTelemetryGetRequest) (respData []responseOp.TrxTelemetryGetResponse, err error) {
	result := repo.getTx(tx).
		Select("id, truck_id, latitude, longitude, temperature, weight_load, vibration, CAST(recorded_at AS VARCHAR) AS recorded_at").
		Table("operations.trx_telemetries")

	if request.TruckId != "" {
		result = result.Where("truck_id = ?", request.TruckId)
	}
	if request.StartDate != "" && request.EndDate != "" {
		result = result.Where("recorded_at >= ? AND recorded_at <= ?", request.StartDate+" 00:00:00", request.EndDate+" 23:59:59")
	}

	result = result.Order("recorded_at DESC")

	if err = result.Scan(&respData).Error; err != nil {
		return []responseOp.TrxTelemetryGetResponse{}, err
	}
	return respData, nil
}

func (repo *TrxTelemetryRepositoryImpl) Save(ctx *gin.Context, tx *gorm.DB, req []entityOp.TrxTelemetry) error {
	result := repo.getTx(tx).Save(&req)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
