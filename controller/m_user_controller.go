package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	requestMst "github.com/royanqodri/fms-ai-hackathon/model/request/master"
	serviceMst "github.com/royanqodri/fms-ai-hackathon/service/master"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type MstUserController interface {
	GetByParams(ctx *gin.Context)
	Post(ctx *gin.Context)
}

type MstUserControllerImpl struct {
	mstUserService serviceMst.MstUserService
}

func NewMstUserController(mstUserService serviceMst.MstUserService) MstUserController {
	return &MstUserControllerImpl{mstUserService: mstUserService}
}

func (controller *MstUserControllerImpl) GetByParams(ctx *gin.Context) {
	req := requestMst.MstUserGetRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	resp, err := controller.mstUserService.GetByParams(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}

func (controller *MstUserControllerImpl) Post(ctx *gin.Context) {
	req := requestMst.MstUserPostRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	err := controller.mstUserService.Post(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, []any{}, "success", 0, 0)
}
