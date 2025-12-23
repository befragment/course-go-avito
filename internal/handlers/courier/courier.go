package courier

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"courier-service/internal/model"
	"courier-service/internal/handlers/utils"
	usecase "courier-service/internal/usecase/courier"
)

type CourierController struct {
	useCase courierUseCase
}

func NewCourierController(useCase courierUseCase) *CourierController {
	return &CourierController{useCase: useCase}
}

func (c *CourierController) GetCourierById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	courier, err := c.useCase.GetCourierById(ctx, id)
	if err != nil {
		if errors.Is(err, usecase.ErrCourierNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, ErrCourierNotFound)
			return
		}
		utils.RespondInternalServerError(w, err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, courier)
}

func (c *CourierController) GetAllCouriers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	couriers, err := c.useCase.GetAllCouriers(ctx)
	if err != nil {
		utils.RespondInternalServerError(w, err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, couriers)
}

func (c *CourierController) CreateCourier(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CourierCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := c.useCase.CreateCourier(ctx, model.Courier{
		Name:          req.Name,
		Phone:         req.Phone,
		TransportType: model.CourierTransportType(req.TransportType),
		Status:        model.CourierStatus(req.Status),
	})

	if err != nil {
		handleCreateError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"id":      strconv.FormatInt(id, 10),
		"message": "Courier created successfully",
	})
}

func (c *CourierController) UpdateCourier(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CourierUpdateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.ID == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, ErrIDRequired)
		return
	}

	courier := req.ToModel()
	err := c.useCase.UpdateCourier(ctx, courier)
	if err != nil {
		handleUpdateError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Courier updated successfully",
	})
}
