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

type MstTruckService interface {
	GetByParams(ctx *gin.Context, request requestMst.MstTruckGetRequest) (respData []responseMst.MstTruckGetResponse, err error)
	Post(ctx *gin.Context, request requestMst.MstTruckPostRequest) (err error)
}

type MstTruckServiceImpl struct {
	mstTruckRepo repoMst.MstTruckRepository
}

func NewMstTruckService(mstTruckRepo repoMst.MstTruckRepository) MstTruckService {
	return &MstTruckServiceImpl{
		mstTruckRepo: mstTruckRepo,
	}
}

func (service MstTruckServiceImpl) GetByParams(ctx *gin.Context, request requestMst.MstTruckGetRequest) (respData []responseMst.MstTruckGetResponse, err error) {
	respData, err = service.mstTruckRepo.GetByParams(ctx, nil, request)
	if err != nil {
		return []responseMst.MstTruckGetResponse{}, fmt.Errorf("failed to get data: %w", err)
	}
	return respData, nil
}

func (service MstTruckServiceImpl) Post(ctx *gin.Context, request requestMst.MstTruckPostRequest) (err error) {
	tx := dbConn.DBConn.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	dataReq := make([]entityMst.MstTruck, len(request.Data))
	for i, data := range request.Data {
		newID := util.ParseUUID(data.Id)
		if newID == uuid.Nil {
			newID = uuid.New()
		}

		dataReq[i] = entityMst.MstTruck{
			Id:          newID,
			CustomerId:  util.MustParseUUID(data.CustomerId),
			PlateNumber: data.PlateNumber,
			CargoType:   data.CargoType,
			Status:      data.Status,
		}
	}

	if len(dataReq) > 0 {
		if err = service.mstTruckRepo.Save(ctx, tx, dataReq); err != nil {
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
