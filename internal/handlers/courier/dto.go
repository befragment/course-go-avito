package courier

import "courier-service/internal/model"

type CourierCreateRequestDTO struct {
	Name          string `json:"name"`
	TransportType string `json:"transport_type"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
}

type CourierUpdateRequestDTO struct {
	ID            int64   `json:"id"`
	TransportType *string `json:"transport_type"`
	Name          *string `json:"name"`
	Phone         *string `json:"phone"`
	Status        *string `json:"status"`
}

func (req CourierUpdateRequestDTO) ToModel() model.Courier {
	courier := model.Courier{ID: req.ID}
	if req.Name != nil {
		courier.Name = *req.Name
	}
	if req.Phone != nil {
		courier.Phone = *req.Phone
	}
	if req.TransportType != nil {
		courier.TransportType = model.CourierTransportType(*req.TransportType)
	}
	if req.Status != nil {
		courier.Status = model.CourierStatus(*req.Status)
	}
	return courier
}
