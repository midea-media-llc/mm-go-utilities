package utils

import (
	geo "github.com/kellydunn/golang-geo"
)

// Distance calculates distance between 2 points
func Distance(point1Lat, point1Lng, point2Lat, point2Lng float64) float64 {
	point1 := geo.NewPoint(point1Lat, point1Lng)
	point2 := geo.NewPoint(point2Lat, point2Lng)

	distance := point1.GreatCircleDistance(point2)

	return distance
}
