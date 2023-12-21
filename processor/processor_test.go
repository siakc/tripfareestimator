package processor


import "testing"



func TestCalcDistanceSpeed(t *testing.T) {
	s :=SegmentData{
		Id:       1,
		Epoch1:   0,
		Epoch2:   3603,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 0,
		Speed:    0,
		Fare:     0,
	}
	d,u := CalcDistanceSpeed(s)
	dwant := 100.0
	uwant := 100.0
    if d > (dwant + 0.1) ||  d < (dwant - 0.1) {
        t.Errorf("CalcDistanceSpeed(complex) returned distance= %f; want %f", d, dwant)
    }
	if u > (uwant + 0.1) || u < (uwant - 0.1) {
        t.Errorf("CalcDistanceSpeed(complex) returned speed= %f; want %f", u, uwant)
    }

}