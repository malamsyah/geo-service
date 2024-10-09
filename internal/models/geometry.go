package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/malamsyah/geo-service/internal/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Type string

const (
	PointType    Type = "Point"
	PolygonType  Type = "Polygon"
	MultiPolygon Type = "MultiPolygon"
)

type Geometry struct {
	Type                    Type `json:"type"`
	PointCoordinates        [2]float64
	PolygonCoordinates      [][][2]float64
	MultiPolygonCoordinates [][][][2]float64
}

func (g Geometry) IsPoint() bool {
	return g.Type == PointType
}

func (g Geometry) IsPolygon() bool {
	return g.Type == PolygonType
}

func (g Geometry) IsMultiPolygon() bool {
	return g.Type == MultiPolygon
}

func (g Geometry) MarshalJSON() ([]byte, error) {
	var coordinates interface{}

	switch g.Type {
	case PointType:
		coordinates = g.PointCoordinates
	case PolygonType:
		coordinates = g.PolygonCoordinates
	case MultiPolygon:
		coordinates = g.MultiPolygonCoordinates
	default:
		coordinates = nil
	}

	return json.Marshal(struct {
		Type        Type        `json:"type"`
		Coordinates interface{} `json:"coordinates"`
	}{
		Type:        g.Type,
		Coordinates: coordinates,
	})
}

func (g *Geometry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if err := json.Unmarshal(raw["type"], &g.Type); err != nil {
		return err
	}

	if g.IsPoint() {
		if err := json.Unmarshal(raw["coordinates"], &g.PointCoordinates); err != nil {
			return err
		}
	}

	if g.IsPolygon() {
		if err := json.Unmarshal(raw["coordinates"], &g.PolygonCoordinates); err != nil {
			return err
		}
	}

	if g.IsMultiPolygon() {
		if err := json.Unmarshal(raw["coordinates"], &g.MultiPolygonCoordinates); err != nil {
			return err
		}
	}

	return nil
}

func (g Geometry) Validate() error {
	if g.IsPoint() {
		return g.validatePoint()
	}

	if g.IsPolygon() {
		return g.validatePolygon()
	}

	return constants.ErrInvalidGeometryType
}

func (g Geometry) validatePoint() error {
	if g.PointCoordinates[0] < -180 || g.PointCoordinates[0] > 180 || g.PointCoordinates[1] < -90 || g.PointCoordinates[1] > 90 {
		return constants.ErrCoordinatesOutOfRange
	}

	return nil
}

func (g Geometry) validatePolygon() error {
	for _, coords := range g.PolygonCoordinates {
		if len(coords) < 2 {
			return constants.ErrInvalidContours
		}

		if coords[0] != coords[len(coords)-1] {
			return constants.ErrInvalidContours
		}

		for _, coord := range coords {
			if coord[0] < -180 || coord[0] > 180 || coord[1] < -90 || coord[1] > 90 {
				return constants.ErrCoordinatesOutOfRange
			}
		}
	}

	return nil
}

func (g Geometry) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	if g.IsPolygon() {
		return clause.Expr{
			SQL:  "ST_PolygonFromText(?)",
			Vars: []interface{}{fmt.Sprintf("POLYGON(%s)", g.polygonCoordinatesToString())},
		}
	}

	if g.IsPoint() {
		return clause.Expr{
			SQL:  "ST_PointFromText(?)",
			Vars: []interface{}{fmt.Sprintf("POINT(%s %s)", fmt.Sprint(g.PointCoordinates[0]), fmt.Sprint(g.PointCoordinates[1]))},
		}
	}

	return clause.Expr{}
}

func (g Geometry) polygonCoordinatesToString() string {
	coords := make([]string, 0)
	for _, c := range g.PolygonCoordinates {
		var points []string

		for _, p := range c {
			points = append(points, fmt.Sprintf("%s %s", fmt.Sprint(p[0]), fmt.Sprint(p[1])))
		}

		joined := strings.Join(points, ",")
		coords = append(coords, fmt.Sprintf("(%s)", joined))
	}

	return strings.Join(coords, ",")
}

func (g *Geometry) Scan(src interface{}) error {
	if src == nil {
		*g = Geometry{}
		return nil
	}

	var data []byte

	switch src := src.(type) {
	case string:
		data = []byte(src)
	case []byte:
		data = src
	default:
		return constants.ErrUnsupportedScan
	}

	return json.Unmarshal(data, g)
}
