package processor

import (
	"github.com/umahmood/haversine"
)
type SegmentData struct{
	Id int64
	Epoch1, Epoch2	int64
	Lat1, Lon1, Lat2, Lon2	float64
	Distance	float64
	Speed	float64
	Fare	float64

}

func CalcDistanceSpeed(input  SegmentData) (float64, float64){

	xc := haversine.Coord{Lat: input.Lat1, Lon: input.Lon1}
	yc := haversine.Coord{Lat: input.Lat2, Lon: input.Lon2}
	_, km := haversine.Distance(xc, yc)
	timeDelta := input.Epoch2-input.Epoch1
	speed := 60 * 60 * km / float64(timeDelta)
	return  km, speed
	
}
