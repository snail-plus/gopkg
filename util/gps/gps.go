package gps

import "math"

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers
	dLat := deg2rad(lat2 - lat1)
	dLon := deg2rad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

// Convert degrees to radians
func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
