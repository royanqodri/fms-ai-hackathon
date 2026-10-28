package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	requestOp "github.com/royanqodri/fms-ai-hackathon/model/request/operation"
	serviceOp "github.com/royanqodri/fms-ai-hackathon/service/operation"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type LogAlertController interface {
	GetByParams(ctx *gin.Context)
	Post(ctx *gin.Context)
}

type LogAlertControllerImpl struct {
	logAlertService serviceOp.LogAlertService
}

func NewLogAlertController(logAlertService serviceOp.LogAlertService) LogAlertController {
	return &LogAlertControllerImpl{logAlertService: logAlertService}
}

func (controller *LogAlertControllerImpl) GetByParams(ctx *gin.Context) {
	req := requestOp.LogAlertGetRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	resp, err := controller.logAlertService.GetByParams(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}

func (controller *LogAlertControllerImpl) Post(ctx *gin.Context) {
	req := requestOp.LogAlertPostRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	err := controller.logAlertService.Post(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, []any{}, "success", 0, 0)
}
