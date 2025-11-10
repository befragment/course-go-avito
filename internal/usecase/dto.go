package usecase

import "courier-service/internal/model"

func courierDBToCourier(courierDB model.CourierDB) model.Courier {
	return model.Courier{
		ID:        courierDB.ID,
		Name:      courierDB.Name,
		Phone:     courierDB.Phone,
		Status:    courierDB.Status,
		CreatedAt: courierDB.CreatedAt,
		UpdatedAt: courierDB.UpdatedAt,
	}
}

func courierUpdateRequestToCourierDB(req model.CourierUpdateRequest) model.CourierDB {
	out := model.CourierDB{ID: req.ID}
	if req.Name != nil   { out.Name = *req.Name }
	if req.Phone != nil  { out.Phone = *req.Phone }
	if req.Status != nil { out.Status = *req.Status }
	return out
}

func courierCreateRequestToCourierDB(req model.CourierCreateRequest) model.CourierDB {
	return model.CourierDB{Name: req.Name, Phone: req.Phone, Status: req.Status}
}