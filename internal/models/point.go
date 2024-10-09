package models

type Point struct {
	ID   uint     `json:"id" gorm:"primaryKey"`
	Data Geometry `json:"data" gorm:"column:data;type:geometry(POINT,4326)"`
}
