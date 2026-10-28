package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	requestMst "github.com/royanqodri/fms-ai-hackathon/model/request/master"
	serviceMst "github.com/royanqodri/fms-ai-hackathon/service/master"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type MstTruckController interface {
	GetByParams(ctx *gin.Context)
	Post(ctx *gin.Context)
}

type MstTruckControllerImpl struct {
	mstTruckService serviceMst.MstTruckService
}

func NewMstTruckController(mstTruckService serviceMst.MstTruckService) MstTruckController {
	return &MstTruckControllerImpl{mstTruckService: mstTruckService}
}

func (controller *MstTruckControllerImpl) GetByParams(ctx *gin.Context) {
	req := requestMst.MstTruckGetRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	resp, err := controller.mstTruckService.GetByParams(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}

func (controller *MstTruckControllerImpl) Post(ctx *gin.Context) {
	req := requestMst.MstTruckPostRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	err := controller.mstTruckService.Post(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, []any{}, "success", 0, 0)
}
