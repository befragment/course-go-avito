package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"courier-service/internal/model"
	"courier-service/internal/usecase"
)

type CourierController struct {
	useCase CourierUseCase
}

func NewCourierController(useCase CourierUseCase) *CourierController {
	return &CourierController{useCase: useCase}
}

func (c *CourierController) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		return
	}

	courier, err := c.useCase.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, usecase.ErrCourierNotFound) {
			http.Error(w, `{"error": "courier not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "` + err.Error() + `"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courier)
}

func (c *CourierController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	couriers, err := c.useCase.GetAll(ctx)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, `{"error": "` + err.Error() + `"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(couriers)
}

func (c *CourierController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var courier model.Courier
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&courier); err != nil {
		http.Error(w, `{"error": "` + err.Error() + `"}`, http.StatusBadRequest)
		return
	}

	id, err := c.useCase.Create(ctx, &model.CourierCreateRequest{
		Name: courier.Name,
		Phone: courier.Phone,
		Status: courier.Status,
	})
	if err != nil {
		switch err {
		case usecase.ErrInvalidCreate: 
			http.Error(w, `{"error": "Missing required fields"}`, http.StatusBadRequest)
		case usecase.ErrInvalidPhoneNumber:
			http.Error(w, `{"error": "Invalid phone number"}`, http.StatusBadRequest)
		case usecase.ErrPhoneNumberExists:
			http.Error(w, `{"error": "Phone number already exists"}`, http.StatusConflict)
		default:
			log.Printf("Internal server error: %v\n", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}
	
	response := map[string]interface{}{
		"id":      id,
		"message": "Profile created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *CourierController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var courier model.Courier
	if err := json.NewDecoder(r.Body).Decode(&courier); err != nil {
		http.Error(w, `{"error": "` + err.Error() + `"}`, http.StatusBadRequest)
		return
	}
	if courier.ID == 0 {
		http.Error(w, `{"error": "id is required"}`, http.StatusBadRequest)
		return
	}
	
	err := c.useCase.Update(ctx, &model.CourierUpdateRequest{
		ID: courier.ID,
		Name: &courier.Name,
		Phone: &courier.Phone,
		Status: &courier.Status,
	})
	if err != nil {
		switch err {
		case usecase.ErrInvalidUpdate:
			http.Error(w, `{"error": "Missing required fields"}`, http.StatusBadRequest)
		case usecase.ErrInvalidPhoneNumber:
			http.Error(w, `{"error": "Invalid phone number"}`, http.StatusBadRequest)
		case usecase.ErrPhoneNumberExists:
			http.Error(w, `{"error": "Phone number already exists"}`, http.StatusConflict)
		default:
			log.Printf("Internal server error: %v\n", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"message": "Courier updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// func RegisterCouriersRoutes(r chi.Router) {
// 	r.Get("/courier/{id}", c.GetById)
// 	r.Get("/couriers", c.GetAll)
// 	r.Post("/courier", c.Create)
// 	r.Put("/courier", c.Update)
// }
