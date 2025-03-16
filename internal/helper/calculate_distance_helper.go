package helper

import "math"

const two = 2
const aHundredEighty = 180

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371

	dLat := (lat2 - lat1) * (math.Pi / aHundredEighty)
	dLon := (lon2 - lon1) * (math.Pi / aHundredEighty)

	a := math.Sin(dLat/two)*math.Sin(dLat/two) +
		math.Cos(lat1*(math.Pi/aHundredEighty))*math.Cos(lat2*(math.Pi/aHundredEighty))*
			math.Sin(dLon/two)*math.Sin(dLon/two)

	c := two * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
