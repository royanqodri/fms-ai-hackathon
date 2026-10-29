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

type MstCustomerService interface {
	GetByParams(ctx *gin.Context, request requestMst.MstCustomerGetRequest) (respData []responseMst.MstCustomerGetResponse, err error)
	Post(ctx *gin.Context, request requestMst.MstCustomerPostRequest) (err error)
}

type MstCustomerServiceImpl struct {
	mstCustomerRepo repoMst.MstCustomerRepository
}

func NewMstCustomerService(mstCustomerRepo repoMst.MstCustomerRepository) MstCustomerService {
	return &MstCustomerServiceImpl{
		mstCustomerRepo: mstCustomerRepo,
	}
}

func (service MstCustomerServiceImpl) GetByParams(ctx *gin.Context, request requestMst.MstCustomerGetRequest) (respData []responseMst.MstCustomerGetResponse, err error) {
	respData, err = service.mstCustomerRepo.GetByParams(ctx, nil, request)
	if err != nil {
		return []responseMst.MstCustomerGetResponse{}, fmt.Errorf("failed to get data: %w", err)
	}
	return respData, nil
}

func (service MstCustomerServiceImpl) Post(ctx *gin.Context, request requestMst.MstCustomerPostRequest) (err error) {
	tx := dbConn.DBConn.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	dataReq := make([]entityMst.MstCustomer, len(request.Data))
	for i, data := range request.Data {
		newID := util.ParseUUID(data.Id)
		if newID == uuid.Nil {
			newID = uuid.New()
		}

		dataReq[i] = entityMst.MstCustomer{
			Id:         newID,
			CustomerNo: data.CustomerNo,
			Name:       data.Name,
			Type:       data.Type,
			PhoneNo:    data.PhoneNo,
		}
	}

	if len(dataReq) > 0 {
		if err = service.mstCustomerRepo.Save(ctx, tx, dataReq); err != nil {
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
