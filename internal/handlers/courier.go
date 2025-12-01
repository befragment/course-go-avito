package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"courier-service/internal/model"
	"courier-service/internal/usecase"
)

type CourierController struct {
	useCase сourierUseCase
}

func NewCourierController(useCase сourierUseCase) *CourierController {
	return &CourierController{useCase: useCase}
}

func (c *CourierController) GetCourierById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	courier, err := c.useCase.GetCourierById(ctx, id)
	if err != nil {
		if errors.Is(err, usecase.ErrCourierNotFound) {
			respondWithError(w, http.StatusNotFound, ErrCourierNotFound)
			return
		}
		respondInternalServerError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, courier)
}

func (c *CourierController) GetAllCouriers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	couriers, err := c.useCase.GetAllCouriers(ctx)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, couriers)
}

func (c *CourierController) CreateCourier(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req model.CourierCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := c.useCase.CreateCourier(ctx, &req)
	if err != nil {
		handleCreateError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"id":      strconv.FormatInt(id, 10),
		"message": "Courier created successfully",
	})
}

func (c *CourierController) UpdateCourier(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req model.CourierUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.ID == 0 {
		respondWithError(w, http.StatusBadRequest, ErrIDRequired)
		return
	}

	err := c.useCase.UpdateCourier(ctx, &req)
	if err != nil {
		handleUpdateError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Courier updated successfully",
	})
}
