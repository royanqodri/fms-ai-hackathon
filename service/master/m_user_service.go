package master

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbConn "github.com/royanqodri/fms-ai-hackathon/database"
	entityMst "github.com/royanqodri/fms-ai-hackathon/model/entity/master"
	requestMst "github.com/royanqodri/fms-ai-hackathon/model/request/master"
	responseMst "github.com/royanqodri/fms-ai-hackathon/model/response/master"
	repoMst "github.com/royanqodri/fms-ai-hackathon/repository/master"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type MstUserService interface {
	GetByParams(ctx *gin.Context, request requestMst.MstUserGetRequest) (respData []responseMst.MstUserGetResponse, err error)
	Post(ctx *gin.Context, request requestMst.MstUserPostRequest) (err error)
}

type MstUserServiceImpl struct {
	mstUserRepo repoMst.MstUserRepository
}

func NewMstUserService(mstUserRepo repoMst.MstUserRepository) MstUserService {
	return &MstUserServiceImpl{
		mstUserRepo: mstUserRepo,
	}
}

func (service MstUserServiceImpl) GetByParams(ctx *gin.Context, request requestMst.MstUserGetRequest) (respData []responseMst.MstUserGetResponse, err error) {
	respData, err = service.mstUserRepo.GetByParams(ctx, nil, request)
	if err != nil {
		return []responseMst.MstUserGetResponse{}, fmt.Errorf("failed to get data: %w", err)
	}
	return respData, nil
}

func (service MstUserServiceImpl) Post(ctx *gin.Context, request requestMst.MstUserPostRequest) (err error) {
	tx := dbConn.DBConn.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	dataReq := make([]entityMst.MstUser, len(request.Data))
	for i, data := range request.Data {
		newID := util.ParseUUID(data.Id)
		if newID == uuid.Nil {
			newID = uuid.New()
		}
		dataReq[i] = entityMst.MstUser{
			Id:         newID,
			CustomerId: util.MustParseUUID(data.CustomerId),
			TelegramId: data.TelegramId,
			Name:       data.Name,
			Role:       data.Role,
		}
	}

	if len(dataReq) > 0 {
		if err = service.mstUserRepo.Save(ctx, tx, dataReq); err != nil {
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
