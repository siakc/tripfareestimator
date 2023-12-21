package processor


import "testing"
import "sync"
import "beat.com/takehome/fareestimator/inputfeeder"

func TestClean(t *testing.T) {

	var wgClean sync.WaitGroup
	inputDatarow := make( chan inputfeeder.Datarow, 10)
	out := make( chan SegmentData, 10)
	want := SegmentData{
		Id:       1,
		Epoch1:   0,
		Epoch2:   3603,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 100,
		Speed:    100,
		Fare:     0,
	}
	go Clean(inputDatarow, out, &wgClean)
	inputDatarow <- inputfeeder.Datarow{
		Id:    1,
		Lat:   52.78836613863794,
		Lon:   8.15902948538737,
		Epoch: 0,
	}
	inputDatarow <- inputfeeder.Datarow{
		Id:    1,
		Lat:   52.88496259502305,
		Lon:   6.677907250263011,
		Epoch: 3603,
	}
	res := <-out


    if res.Distance > (want.Distance + 0.1) ||  res.Distance < (want.Distance - 0.1) {
        t.Errorf("Clean(complex) returned distance= %f; want %f", res.Distance, want.Distance)
    }
	if res.Speed > (want.Speed + 0.1) || res.Speed  < (want.Speed - 0.1) {
        t.Errorf("Clean(complex) returned speed= %f; want %f", res.Speed, want.Speed)
    }

	inputDatarow <- inputfeeder.Datarow{
		Id:    2,
		Lat:   52.78836613863794,
		Lon:   8.15902948538737,
		Epoch: 0,
	}
	inputDatarow <- inputfeeder.Datarow{
		Id:    2,
		Lat:   52.88496259502305,
		Lon:   16.677907250263011,
		Epoch: 3603,
	}
	inputDatarow <- inputfeeder.Datarow{
		Id:    2,
		Lat:   52.88496259502305,
		Lon:   6.677907250263011,
		Epoch: 3603,
	}
	res = <-out


    if res.Distance > (want.Distance + 0.1) ||  res.Distance < (want.Distance - 0.1) {
        t.Errorf("Clean(complex) returned distance= %f; want %f", res.Distance, want.Distance)
    }
	if res.Speed > (want.Speed + 0.1) || res.Speed  < (want.Speed - 0.1) {
        t.Errorf("Clean(complex) returned speed= %f; want %f", res.Speed, want.Speed)
    }

}

func TestPrepareAndClean(t *testing.T) {


	inputDatarow := make( chan inputfeeder.Datarow, 10)
	out := make( chan SegmentData, 10)
	go PrepareAndClean(inputDatarow, out)
	inputDatarow <- inputfeeder.Datarow{
		Id:    1,
		Lat:   52.78836613863794,
		Lon:   8.15902948538737,
		Epoch: 0,
	}
	inputDatarow <- inputfeeder.Datarow{
		Id:    1,
		Lat:   52.88496259502305,
		Lon:   6.677907250263011,
		Epoch: 3603,
	}
	res := <-out
	want := SegmentData{
		Id:       1,
		Epoch1:   0,
		Epoch2:   3603,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 100,
		Speed:    100,
		Fare:     0,
	}

    if res.Distance > (want.Distance + 0.1) ||  res.Distance < (want.Distance - 0.1) {
        t.Errorf("Clean(complex) returned distance= %f; want %f", res.Distance, want.Distance)
    }
	if res.Speed > (want.Speed + 0.1) || res.Speed  < (want.Speed - 0.1) {
        t.Errorf("Clean(complex) returned speed= %f; want %f", res.Speed, want.Speed)
    }

}