package operation

import (
	"fmt"

	"github.com/gin-gonic/gin"

	dbConn "github.com/royanqodri/fms-ai-hackathon/database"
	entityOp "github.com/royanqodri/fms-ai-hackathon/model/entity/operation"
	requestOp "github.com/royanqodri/fms-ai-hackathon/model/request/operation"
	responseOp "github.com/royanqodri/fms-ai-hackathon/model/response/operation"
	repoOp "github.com/royanqodri/fms-ai-hackathon/repository/operation"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type TrxTelemetryService interface {
	GetByParams(ctx *gin.Context, request requestOp.TrxTelemetryGetRequest) (respData []responseOp.TrxTelemetryGetResponse, err error)
	Post(ctx *gin.Context, request requestOp.TrxTelemetryPostRequest) (err error)
}

type TrxTelemetryServiceImpl struct {
	trxTelemetryRepo repoOp.TrxTelemetryRepository
}

func NewTrxTelemetryService(trxTelemetryRepo repoOp.TrxTelemetryRepository) TrxTelemetryService {
	return &TrxTelemetryServiceImpl{
		trxTelemetryRepo: trxTelemetryRepo,
	}
}

func (service TrxTelemetryServiceImpl) GetByParams(ctx *gin.Context, request requestOp.TrxTelemetryGetRequest) (respData []responseOp.TrxTelemetryGetResponse, err error) {
	respData, err = service.trxTelemetryRepo.GetByParams(ctx, nil, request)
	if err != nil {
		return []responseOp.TrxTelemetryGetResponse{}, fmt.Errorf("failed to get data: %w", err)
	}
	return respData, nil
}

func (service TrxTelemetryServiceImpl) Post(ctx *gin.Context, request requestOp.TrxTelemetryPostRequest) (err error) {
	tx := dbConn.DBConn.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	dataReq := make([]entityOp.TrxTelemetry, len(request.Data))
	for i, data := range request.Data {
		dataReq[i] = entityOp.TrxTelemetry{
			Id:          util.ParseUUID(data.Id),
			TruckId:     util.MustParseUUID(data.TruckId),
			Latitude:    data.Latitude,
			Longitude:   data.Longitude,
			Temperature: data.Temperature,
			WeightLoad:  data.WeightLoad,
			Vibration:   data.Vibration,
		}
	}

	if len(dataReq) > 0 {
		if err = service.trxTelemetryRepo.Save(ctx, tx, dataReq); err != nil {
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	// TODO: GET Meta AI
	return nil
}
