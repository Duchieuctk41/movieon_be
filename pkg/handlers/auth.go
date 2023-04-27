package handlers

import (
	"github.com/praslar/cloud0/ginext"
	"github.com/praslar/cloud0/logger"
	"movieon_be/pkg/model"
	"movieon_be/pkg/utils"
	"net/http"
)

// Login
// @Tags Login
// @Accept  json
// @Produce  json
// @Param data body model.BlacklistParam true "body data"
// @Success 200 {object} interface{}
// @Router /api/v1/auth/login [post]
func (h *UserHandlers) Login(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))

	req := model.LoginRequest{}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.WithError(err).Error("error_400: error parse")
		return nil, ginext.NewError(http.StatusUnauthorized, "Yêu cầu không hợp lệ")
	}

	data, err := h.service.Login(r.Context(), req)
	if err != nil {
		return nil, ginext.NewError(http.StatusUnauthorized, "Tài khoản hoặc mật khẩu không chính xác!")
	}
	return ginext.NewResponseData(http.StatusOK, data), nil
}

// Register
// @Tags Login
// @Accept  json
// @Produce  json
// @Param data body model.BlacklistParam true "body data"
// @Success 200 {object} interface{}
// @Router /api/v1/auth/register [post]
func (h *UserHandlers) Register(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))

	req := model.RegisterRequest{}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.WithError(err).Error("error_400: error parse")
		return nil, ginext.NewError(http.StatusUnauthorized, "Yêu cầu không hợp lệ")
	}

	if err := h.service.Register(r.Context(), req); err != nil {
		return nil, err
	}
	return ginext.NewResponse(http.StatusOK), nil
}
