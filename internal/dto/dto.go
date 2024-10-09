package dto

import "github.com/malamsyah/geo-service/internal/models"

type Response struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  interface{} `json:"results,omitempty"`
}

type CreatePointRequest struct {
	Data models.Geometry `json:"data" binding:"required"`
}

func (r CreatePointRequest) ToModel() models.Point {
	return models.Point{Data: r.Data}
}

type CreateContourRequest struct {
	Data models.Geometry `json:"data" binding:"required"`
}

func (r CreateContourRequest) ToModel() models.Contour {
	return models.Contour{Data: r.Data}
}
