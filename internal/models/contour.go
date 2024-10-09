package models

type Contour struct {
	ID   uint     `json:"id,omitempty" gorm:"primaryKey"`
	Data Geometry `json:"data" gorm:"column:data;type:geometry(POLYGON,4326)"`
}
