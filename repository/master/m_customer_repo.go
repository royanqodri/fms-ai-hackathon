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

type MstCustomerRepository interface {
	GetById(ctx *gin.Context, tx *gorm.DB, id string) (responseMst.MstCustomerGetResponse, error)
	GetByParams(ctx *gin.Context, tx *gorm.DB, request requestMst.MstCustomerGetRequest) (respData []responseMst.MstCustomerGetResponse, err error)
	Save(ctx *gin.Context, tx *gorm.DB, req []entityMst.MstCustomer) error
}

type MstCustomerRepositoryImpl struct{}

func NewMstCustomerRepository() MstCustomerRepository {
	return &MstCustomerRepositoryImpl{}
}

func (repo MstCustomerRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return dbConn.DBConn
	}
	return tx
}

func (repo MstCustomerRepositoryImpl) GetById(ctx *gin.Context, tx *gorm.DB, id string) (responseMst.MstCustomerGetResponse, error) {
	var respData responseMst.MstCustomerGetResponse

	result := repo.getTx(tx).
		Select("id, customer_no, name, type, phone_no, CAST(created_at AS VARCHAR) AS created_at, CAST(updated_at AS VARCHAR) AS updated_at").
		Table("master.mst_customers").
		Where("id = ?", id).
		Scan(&respData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return responseMst.MstCustomerGetResponse{}, nil
		}
		return responseMst.MstCustomerGetResponse{}, result.Error
	}
	return respData, nil
}

func (repo MstCustomerRepositoryImpl) GetByParams(ctx *gin.Context, tx *gorm.DB, request requestMst.MstCustomerGetRequest) (respData []responseMst.MstCustomerGetResponse, err error) {
	result := repo.getTx(tx).
		Select("id, customer_no, name, type, phone_no, CAST(created_at AS VARCHAR) AS created_at, CAST(updated_at AS VARCHAR) AS updated_at").
		Table("master.mst_customers")

	if request.Search != "" {
		result = result.Where("name ILIKE ? OR customer_no ILIKE ?", "%"+request.Search+"%", "%"+request.Search+"%")
	}
	if request.Type != "" {
		result = result.Where("type = ?", request.Type)
	}

	result = result.Order("created_at DESC")

	if err = result.Scan(&respData).Error; err != nil {
		return []responseMst.MstCustomerGetResponse{}, err
	}
	return respData, nil
}

func (repo *MstCustomerRepositoryImpl) Save(ctx *gin.Context, tx *gorm.DB, req []entityMst.MstCustomer) error {
	result := repo.getTx(tx).Create(&req)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
