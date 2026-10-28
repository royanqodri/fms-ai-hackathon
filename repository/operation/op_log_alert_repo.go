package operation

import (
	"errors"

	"github.com/gin-gonic/gin"
	dbConn "github.com/royanqodri/fms-ai-hackathon/database"
	"gorm.io/gorm"

	entityOp "github.com/royanqodri/fms-ai-hackathon/model/entity/operation"
	requestOp "github.com/royanqodri/fms-ai-hackathon/model/request/operation"
	responseOp "github.com/royanqodri/fms-ai-hackathon/model/response/operation"
)

type LogAlertRepository interface {
	GetById(ctx *gin.Context, tx *gorm.DB, id string) (responseOp.LogAlertGetResponse, error)
	GetByParams(ctx *gin.Context, tx *gorm.DB, request requestOp.LogAlertGetRequest) (respData []responseOp.LogAlertGetResponse, err error)
	Save(ctx *gin.Context, tx *gorm.DB, req []entityOp.LogAlert) error
	UpdateResolvedStatus(ctx *gin.Context, tx *gorm.DB, id string, isResolved bool, updateBy string) error
}

type LogAlertRepositoryImpl struct{}

func NewLogAlertRepository() LogAlertRepository {
	return &LogAlertRepositoryImpl{}
}

func (repo LogAlertRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return dbConn.DBConn
	}
	return tx
}

func (repo LogAlertRepositoryImpl) GetById(ctx *gin.Context, tx *gorm.DB, id string) (responseOp.LogAlertGetResponse, error) {
	var respData responseOp.LogAlertGetResponse

	result := repo.getTx(tx).
		Select("id, truck_id, alert_type, ai_analysis, is_resolved, CAST(created_at AS VARCHAR) AS created_at").
		Table("operations.log_alerts").
		Where("id = ?", id).
		Scan(&respData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return responseOp.LogAlertGetResponse{}, nil
		}
		return responseOp.LogAlertGetResponse{}, result.Error
	}
	return respData, nil
}

func (repo LogAlertRepositoryImpl) GetByParams(ctx *gin.Context, tx *gorm.DB, request requestOp.LogAlertGetRequest) (respData []responseOp.LogAlertGetResponse, err error) {
	result := repo.getTx(tx).
		Select("id, truck_id, alert_type, ai_analysis, is_resolved, CAST(created_at AS VARCHAR) AS created_at").
		Table("operations.log_alerts")

	if request.TruckId != "" {
		result = result.Where("truck_id = ?", request.TruckId)
	}
	if request.AlertType != "" {
		result = result.Where("alert_type = ?", request.AlertType)
	}
	if request.IsResolved != nil {
		result = result.Where("is_resolved = ?", *request.IsResolved)
	}

	result = result.Order("created_at DESC")

	if err = result.Scan(&respData).Error; err != nil {
		return []responseOp.LogAlertGetResponse{}, err
	}
	return respData, nil
}

func (repo *LogAlertRepositoryImpl) Save(ctx *gin.Context, tx *gorm.DB, req []entityOp.LogAlert) error {
	result := repo.getTx(tx).Save(&req)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *LogAlertRepositoryImpl) UpdateResolvedStatus(ctx *gin.Context, tx *gorm.DB, id string, isResolved bool, updateBy string) error {
	result := repo.getTx(tx).
		Table("operations.log_alerts").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_resolved": isResolved,
			"updated_at":  gorm.Expr("NOW()"),
		})

	if result.Error != nil {
		return result.Error
	}
	return nil
}
