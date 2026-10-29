package operation

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbConn "github.com/royanqodri/fms-ai-hackathon/database"
	entityOp "github.com/royanqodri/fms-ai-hackathon/model/entity/operation"
	requestOp "github.com/royanqodri/fms-ai-hackathon/model/request/operation"
	responseOp "github.com/royanqodri/fms-ai-hackathon/model/response/operation"
	repoOp "github.com/royanqodri/fms-ai-hackathon/repository/operation"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type LogAlertService interface {
	GetByParams(ctx *gin.Context, request requestOp.LogAlertGetRequest) (respData []responseOp.LogAlertGetResponse, err error)
	Post(ctx *gin.Context, request requestOp.LogAlertPostRequest) (err error)
	UpdateResolved(ctx *gin.Context, id string, isResolved bool, updateBy string) (err error)
}

type LogAlertServiceImpl struct {
	logAlertRepo repoOp.LogAlertRepository
}

func NewLogAlertService(logAlertRepo repoOp.LogAlertRepository) LogAlertService {
	return &LogAlertServiceImpl{
		logAlertRepo: logAlertRepo,
	}
}

func (service LogAlertServiceImpl) GetByParams(ctx *gin.Context, request requestOp.LogAlertGetRequest) (respData []responseOp.LogAlertGetResponse, err error) {
	respData, err = service.logAlertRepo.GetByParams(ctx, nil, request)
	if err != nil {
		return []responseOp.LogAlertGetResponse{}, fmt.Errorf("failed to get data: %w", err)
	}
	return respData, nil
}

func (service LogAlertServiceImpl) Post(ctx *gin.Context, request requestOp.LogAlertPostRequest) (err error) {
	tx := dbConn.DBConn.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	dataReq := make([]entityOp.LogAlert, len(request.Data))
	for i, data := range request.Data {
		newID := util.ParseUUID(data.Id)
		if newID == uuid.Nil {
			newID = uuid.New()
		}
		dataReq[i] = entityOp.LogAlert{
			Id:         newID,
			TruckId:    util.MustParseUUID(data.TruckId),
			AlertType:  data.AlertType,
			AiAnalysis: data.AiAnalysis,
			IsResolved: data.IsResolved,
		}
	}

	if len(dataReq) > 0 {
		if err = service.logAlertRepo.Save(ctx, tx, dataReq); err != nil {
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (service LogAlertServiceImpl) UpdateResolved(ctx *gin.Context, id string, isResolved bool, updateBy string) (err error) {
	tx := dbConn.DBConn.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = service.logAlertRepo.UpdateResolvedStatus(ctx, tx, id, isResolved, updateBy); err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
