package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	requestOp "github.com/royanqodri/fms-ai-hackathon/model/request/operation"
	serviceOp "github.com/royanqodri/fms-ai-hackathon/service/operation"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type TrxTelemetryController interface {
	GetByParams(ctx *gin.Context)
	Post(ctx *gin.Context)
}

type TrxTelemetryControllerImpl struct {
	trxTelemetryService serviceOp.TrxTelemetryService
}

func NewTrxTelemetryController(trxTelemetryService serviceOp.TrxTelemetryService) TrxTelemetryController {
	return &TrxTelemetryControllerImpl{trxTelemetryService: trxTelemetryService}
}

func (controller *TrxTelemetryControllerImpl) GetByParams(ctx *gin.Context) {
	req := requestOp.TrxTelemetryGetRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	resp, err := controller.trxTelemetryService.GetByParams(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}

func (controller *TrxTelemetryControllerImpl) Post(ctx *gin.Context) {
	req := requestOp.TrxTelemetryPostRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	err := controller.trxTelemetryService.Post(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, []any{}, "success", 0, 0)
}
