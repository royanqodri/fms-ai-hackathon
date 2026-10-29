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

type MstUserRepository interface {
	GetById(ctx *gin.Context, tx *gorm.DB, id string) (responseMst.MstUserGetResponse, error)
	GetByParams(ctx *gin.Context, tx *gorm.DB, request requestMst.MstUserGetRequest) (respData []responseMst.MstUserGetResponse, err error)
	GetByTelegramId(ctx *gin.Context, tx *gorm.DB, telegramId string) (entityMst.MstUser, error)
	Save(ctx *gin.Context, tx *gorm.DB, req []entityMst.MstUser) error
}

type MstUserRepositoryImpl struct{}

func NewMstUserRepository() MstUserRepository {
	return &MstUserRepositoryImpl{}
}

func (repo MstUserRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return dbConn.DBConn
	}
	return tx
}

func (repo MstUserRepositoryImpl) GetById(ctx *gin.Context, tx *gorm.DB, id string) (responseMst.MstUserGetResponse, error) {
	var respData responseMst.MstUserGetResponse

	result := repo.getTx(tx).
		Select("id, customer_id, telegram_id, name, role, CAST(created_at AS VARCHAR) AS created_at, CAST(updated_at AS VARCHAR) AS updated_at").
		Table("master.mst_users").
		Where("id = ?", id).
		Scan(&respData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return responseMst.MstUserGetResponse{}, nil
		}
		return responseMst.MstUserGetResponse{}, result.Error
	}
	return respData, nil
}

func (repo MstUserRepositoryImpl) GetByParams(ctx *gin.Context, tx *gorm.DB, request requestMst.MstUserGetRequest) (respData []responseMst.MstUserGetResponse, err error) {
	result := repo.getTx(tx).
		Select("id, customer_id, telegram_id, name, role, CAST(created_at AS VARCHAR) AS created_at, CAST(updated_at AS VARCHAR) AS updated_at").
		Table("master.mst_users")

	if request.Search != "" {
		result = result.Where("name ILIKE ? OR telegram_id ILIKE ?", "%"+request.Search+"%", "%"+request.Search+"%")
	}
	if request.Role != "" {
		result = result.Where("role = ?", request.Role)
	}

	result = result.Order("created_at DESC")

	if err = result.Scan(&respData).Error; err != nil {
		return []responseMst.MstUserGetResponse{}, err
	}
	return respData, nil
}

func (repo MstUserRepositoryImpl) GetByTelegramId(ctx *gin.Context, tx *gorm.DB, telegramId string) (entityMst.MstUser, error) {
	var user entityMst.MstUser

	result := repo.getTx(tx).
		Table("master.mst_users").
		Where("telegram_id = ?", telegramId).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entityMst.MstUser{}, nil
		}
		return entityMst.MstUser{}, result.Error
	}
	return user, nil
}

func (repo *MstUserRepositoryImpl) Save(ctx *gin.Context, tx *gorm.DB, req []entityMst.MstUser) error {
	result := repo.getTx(tx).Create(&req)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
