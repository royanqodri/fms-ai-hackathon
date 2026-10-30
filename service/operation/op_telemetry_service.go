package operation

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
	logAlertRepo     repoOp.LogAlertRepository
}

func NewTrxTelemetryService(
	trxTelemetryRepo repoOp.TrxTelemetryRepository,
	logAlertRepo repoOp.LogAlertRepository) TrxTelemetryService {
	return &TrxTelemetryServiceImpl{
		trxTelemetryRepo: trxTelemetryRepo,
		logAlertRepo:     logAlertRepo,
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
		newID := util.ParseUUID(data.Id)
		if newID == uuid.Nil {
			newID = uuid.New()
		}

		dataReq[i] = entityOp.TrxTelemetry{
			Id:          newID,
			TruckId:     util.MustParseUUID(data.TruckId),
			Latitude:    data.Latitude,
			Longitude:   data.Longitude,
			Temperature: data.Temperature,
			WeightLoad:  data.WeightLoad,
			Vibration:   data.Vibration,
			FuelLevel:   data.FuelLevel,
			OilPressure: data.OilPressure,
			EngineRpm:   data.EngineRpm,
			RecordedAt:  time.Now(),
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

	for _, data := range dataReq {
		if data.Temperature > 38.0 {
			go service.analyzeAndAlert(data)
		}
	}

	return nil
}

func (service TrxTelemetryServiceImpl) analyzeAndAlert(data entityOp.TrxTelemetry) {
	systemPrompt, userPrompt := util.GenerateTelemetryPrompt(
		data.TruckId.String(),
		data.Latitude,
		data.Longitude,
		data.Temperature,
		data.WeightLoad,
		data.Vibration,
		data.FuelLevel,
		data.OilPressure,
		data.EngineRpm,
	)

	// GET AI
	aiReply, err := util.AskMetaAI(systemPrompt, userPrompt)
	if err != nil {
		log.Printf("❌ AI Error: %v", err)
		return
	}

	alertType := "WARNING"
	if data.Temperature > 45.0 {
		alertType = "CRITICAL_OVERHEAT"
	}

	alertData := []entityOp.LogAlert{
		{
			Id:         uuid.New(),
			TruckId:    data.TruckId,
			AlertType:  alertType,
			AiAnalysis: aiReply,
			IsResolved: false,
		},
	}

	txAlert := dbConn.DBConn.Begin()
	if err := service.logAlertRepo.Save(nil, txAlert, alertData); err != nil {
		txAlert.Rollback()
		log.Printf("❌ Failed to save alert: %v", err)
		return
	}
	txAlert.Commit()

	log.Printf("✅ AI Alert Generated for Truck %s: %s", data.TruckId.String(), alertType)
}
