package master

import (
	"errors"

	"github.com/gin-gonic/gin"
	dbConn "github.com/royanqodri/fms-ai-hackathon/database"
	"gorm.io/gorm"

	entityMst "github.com/royanqodri/fms-ai-hackathon/model/entity/master"
	requestMst "github.com/royanqodri/fms-ai-hackathon/model/request/master"
	responseMst "github.com/royanqodri/fms-ai-hackathon/model/response/master"
)

type MstTruckRepository interface {
	GetById(ctx *gin.Context, tx *gorm.DB, id string) (responseMst.MstTruckGetResponse, error)
	GetByParams(ctx *gin.Context, tx *gorm.DB, request requestMst.MstTruckGetRequest) (respData []responseMst.MstTruckGetResponse, err error)
	Save(ctx *gin.Context, tx *gorm.DB, req []entityMst.MstTruck) error
}

type MstTruckRepositoryImpl struct{}

func NewMstTruckRepository() MstTruckRepository {
	return &MstTruckRepositoryImpl{}
}

func (repo MstTruckRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return dbConn.DBConn
	}
	return tx
}

func (repo MstTruckRepositoryImpl) GetById(ctx *gin.Context, tx *gorm.DB, id string) (responseMst.MstTruckGetResponse, error) {
	var respData responseMst.MstTruckGetResponse

	result := repo.getTx(tx).
		Select("id, customer_id, plate_number, cargo_type, status, CAST(created_at AS VARCHAR) AS created_at, CAST(updated_at AS VARCHAR) AS updated_at").
		Table("master.mst_trucks").
		Where("id = ?", id).
		Scan(&respData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return responseMst.MstTruckGetResponse{}, nil
		}
		return responseMst.MstTruckGetResponse{}, result.Error
	}
	return respData, nil
}

func (repo MstTruckRepositoryImpl) GetByParams(ctx *gin.Context, tx *gorm.DB, request requestMst.MstTruckGetRequest) (respData []responseMst.MstTruckGetResponse, err error) {
	result := repo.getTx(tx).
		Select("id, customer_id, plate_number, cargo_type, status, CAST(created_at AS VARCHAR) AS created_at, CAST(updated_at AS VARCHAR) AS updated_at").
		Table("master.mst_trucks")

	if request.Search != "" {
		result = result.Where("plate_number ILIKE ?", "%"+request.Search+"%")
	}
	if request.CargoType != "" {
		result = result.Where("cargo_type = ?", request.CargoType)
	}
	if request.Status != "" {
		result = result.Where("status = ?", request.Status)
	}
	if request.CustomerId != "" {
		result = result.Where("customer_id = ?", request.CustomerId)
	}

	result = result.Order("created_at DESC")

	if err = result.Scan(&respData).Error; err != nil {
		return []responseMst.MstTruckGetResponse{}, err
	}
	return respData, nil
}

func (repo *MstTruckRepositoryImpl) Save(ctx *gin.Context, tx *gorm.DB, req []entityMst.MstTruck) error {
	result := repo.getTx(tx).Save(&req)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
