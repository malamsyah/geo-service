package constants

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInternal = errors.New("internal error")
var ErrInvalidPoint = errors.New("invalid point")
var ErrCoordinatesOutOfRange = errors.New("coordinates out of range")
var ErrInvalidGeometryType = errors.New("invalid geometry type")
var ErrUnsupportedScan = errors.New("unsupported scan")
var ErrInvalidContours = errors.New("invalid contours")
var ErrContourNotFound = errors.New("contour not found")
