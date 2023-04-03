package handlers

import (
	"github.com/praslar/cloud0/ginext"
	"github.com/praslar/cloud0/logger"
	"github.com/praslar/lib/common"
	"movieon_be/pkg/model"
	"movieon_be/pkg/service"
	"movieon_be/pkg/utils"
	"net/http"
)

type MovieHandlers struct {
	service service.MovieInterface
}

func NewMovieHandlers(service service.MovieInterface) *MovieHandlers {
	return &MovieHandlers{service: service}
}

// Create
// @Tags Create
// @Accept  json
// @Produce  json
// @Param data body model.MovieRequest true "body data"
// @Success 200 {object} interface{}
// @Router /api/v1/movie/create [post]
func (h *MovieHandlers) Create(r *ginext.Request) (*ginext.Response, error) {
	req := model.MovieRequest{}
	r.MustBind(&req)

	if err := common.CheckRequireValid(req); err != nil {
		return nil, ginext.NewError(http.StatusBadRequest, utils.MessageError()[http.StatusBadRequest])
	}

	data, err := h.service.Create(r.Context(), req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusOK, data), nil
}

// Update
// @Tags Update
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param data body model.MovieRequest true "body data"
// @Success 200 {object} interface{}
// @Router /api/v1/movie/update/:id [put]
func (h *MovieHandlers) Update(r *ginext.Request) (*ginext.Response, error) {
	id := utils.GetIdFromUri(r.GinCtx)
	if id == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	req := model.MovieRequest{}
	r.MustBind(&req)
	req.ID = id

	if err := common.CheckRequireValid(req); err != nil {
		return nil, ginext.NewError(http.StatusBadRequest, utils.MessageError()[http.StatusBadRequest])
	}

	data, err := h.service.Update(r.Context(), req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusOK, data), nil
}

// Delete
// @Tags Delete
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} interface{}
// @Router /api/v1/movie/delete/:id [delete]
func (h *MovieHandlers) Delete(r *ginext.Request) (*ginext.Response, error) {
	id := utils.ParseIDFromUri(r.GinCtx)
	if id == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	if err := h.service.Delete(r.Context(), id.String()); err != nil {
		return nil, err
	}
	return ginext.NewResponse(http.StatusOK), nil
}

// GetOne
// @Tags GetOne
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} interface{}
// @Router /api/v1/movie/get-one/:id [get]
func (h *MovieHandlers) GetOne(r *ginext.Request) (*ginext.Response, error) {

	id := utils.ParseIDFromUri(r.GinCtx)
	if id == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	data, err := h.service.GetOne(r.Context(), id.String())
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusOK, data), nil
}

// GetList
// @Tags GetList
// @Accept  json
// @Produce  json
// @Param data body model.BlacklistParam true "body data"
// @Success 200 {object} interface{}
// @Router /api/v1/movie/get-list [get]
func (h *MovieHandlers) GetList(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))

	req := model.MovieParams{}
	if err := r.GinCtx.BindQuery(&req); err != nil {
		log.WithError(err).Error("error_400: error parse")
		return nil, ginext.NewError(http.StatusBadRequest, "Yêu cầu không hợp lệ")
	}

	data, err := h.service.GetList(r.Context(), req)
	if err != nil {
		return nil, err
	}
	return &ginext.Response{Code: http.StatusOK, GeneralBody: &ginext.GeneralBody{
		Data: data.Data,
		Meta: data.Meta,
	}}, nil
}
