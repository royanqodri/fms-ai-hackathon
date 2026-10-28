package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	requestMst "github.com/royanqodri/fms-ai-hackathon/model/request/master"
	serviceMst "github.com/royanqodri/fms-ai-hackathon/service/master"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type MstCustomerController interface {
	GetByParams(ctx *gin.Context)
	Post(ctx *gin.Context)
}

type MstCustomerControllerImpl struct {
	mstCustomerService serviceMst.MstCustomerService
}

func NewMstCustomerController(mstCustomerService serviceMst.MstCustomerService) MstCustomerController {
	return &MstCustomerControllerImpl{mstCustomerService: mstCustomerService}
}

func (controller *MstCustomerControllerImpl) GetByParams(ctx *gin.Context) {
	req := requestMst.MstCustomerGetRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	resp, err := controller.mstCustomerService.GetByParams(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}

func (controller *MstCustomerControllerImpl) Post(ctx *gin.Context) {
	req := requestMst.MstCustomerPostRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	err := controller.mstCustomerService.Post(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, []any{}, "success", 0, 0)
}
