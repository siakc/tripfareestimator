package estimator

import "testing"
import "sync"
import "beat.com/takehome/fareestimator/processor"

func TestEstimateSegmentFare(t *testing.T){
	var wg sync.WaitGroup

	input := make( chan processor.SegmentData, 10)
	output := make( chan processor.SegmentData, 10)
	s1 := processor.SegmentData{
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
	input<-s1
	
	wg.Add(1)
	go EstimateSegmentFare(input, output, &wg)
	go func(){
		defer close(output)
		wg.Wait()
	}()
	res := <-output
	want := 130.0
	if res.Fare > (want + 0.1) ||  res.Fare < (want - 0.1) {
        t.Errorf("EstimateSegmentFare(complex) returned fare= %f; want %f", res.Fare, want)
    }
	s2 := processor.SegmentData{
		Id:       1,
		Epoch1:   5 * 60 * 60 + 1,
		Epoch2:   3603 + 5 * 60 * 60 + 1,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 100,
		Speed:    100,
		Fare:     0,
	}

	input<-s2
	res = <-output
	want = 74
	if res.Fare > (want + 0.1) ||  res.Fare < (want - 0.1) {
        t.Errorf("EstimateSegmentFare(complex) returned fare= %f; want %f", res.Fare, want)
    }

	s3 := processor.SegmentData{
		Id:       1,
		Epoch1:   5 * 60 * 60 + 1,
		Epoch2:   3603 + 5 * 60 * 60 + 1,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 100,
		Speed:    5,
		Fare:     0,
	}

	input<-s3
	res = <-output
	want = 11.9
	if res.Fare > (want + 0.1) ||  res.Fare < (want - 0.1) {
        t.Errorf("EstimateSegmentFare(complex) returned fare= %f; want %f", res.Fare, want)
    }

	close(input)
}

func TestMakeEstimates(t *testing.T){
	var wg sync.WaitGroup
	var wgSeg sync.WaitGroup
	input := make( chan processor.SegmentData, 10)
	segOutput := make( chan  processor.SegmentData, 10)
	totalOutput := make( chan Estimate, 10)
	


	s1 := processor.SegmentData{
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

	s2 := processor.SegmentData{
		Id:       2,
		Epoch1:   5 * 60 * 60 + 1,
		Epoch2:   3603 + 5 * 60 * 60 + 1,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 100,
		Speed:    100,
		Fare:     0,
	}
	s3 := processor.SegmentData{
		Id:       3,
		Epoch1:   5 * 60 * 60 + 1,
		Epoch2:   3603 + 5 * 60 * 60 + 1,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 100,
		Speed:    5,
		Fare:     0,
	}

	s4 := processor.SegmentData{
		Id:       4,
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

	s5 := processor.SegmentData{
		Id:       4,
		Epoch1:   5 * 60 * 60 + 1,
		Epoch2:   3603 + 5 * 60 * 60 + 1,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 100,
		Speed:    100,
		Fare:     0,
	}
	s6 := processor.SegmentData{
		Id:       4,
		Epoch1:   5 * 60 * 60 + 1,
		Epoch2:   3603 + 5 * 60 * 60 + 1,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 100,
		Speed:    5,
		Fare:     0,
	}
	s7 := processor.SegmentData{
		Id:       5,
		Epoch1:   5 * 60 * 60 + 1,
		Epoch2:   3603 + 5 * 60 * 60 + 1,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 10,
		Speed:    11,
		Fare:     0,
	}
	s8 := processor.SegmentData{
		Id:       6,
		Epoch1:   5 * 60 * 60 + 1,
		Epoch2:   3603 + 5 * 60 * 60 + 1,
		Lat1:     52.78836613863794,
		Lon1:     8.15902948538737,
		Lat2:     52.88496259502305,
		Lon2:     6.677907250263011,
		Distance: 2,
		Speed:    11,
		Fare:     0,
	}

	wgSeg.Add(1)
	go EstimateSegmentFare(input, segOutput, &wgSeg)
	go func(){
		defer close(segOutput)
		wgSeg.Wait()
	}()
		
	wg.Add(1)
	go MakeEstimates(segOutput, totalOutput, &wg)
	go func(){
		defer close(totalOutput)
		wg.Wait()
	}()

	input<-s1
	input<-s2
	input<-s3
	input<-s4
	input<-s5
	input<-s6
	input<-s7
	input<-s8
	close(input)

	var want float64
	for res := range totalOutput{
		if res.Id == 1 {
			want = 130.0 + 1.3
			if res.Fare > (want + 0.1) ||  res.Fare < (want - 0.1) {
				t.Errorf("MakeEstimates(complex) returned fare= %f; want %f", res.Fare, want)
			}
		} else if res.Id == 2 {
			want = 74 + 1.3
			if res.Fare > (want + 0.1) ||  res.Fare < (want - 0.1) {
				t.Errorf("MakeEstimates(complex) returned fare= %f; want %f", res.Fare, want)
			}
		} else if res.Id == 3 {
			want = 11.9 + 1.3
			if res.Fare > (want + 0.1) ||  res.Fare < (want - 0.1) {
				t.Errorf("MakeEstimates(complex) returned fare= %f; want %f", res.Fare, want)
			}

		} else  if res.Id == 4 {
			want = 130.0 + 74 + 11.9 + 1.3
			if res.Fare > (want + 0.1) ||  res.Fare < (want - 0.1) {
				t.Errorf("MakeEstimates(complex) returned fare= %f; want %f", res.Fare, want)
			}

		}else  if res.Id == 5 {
			want = 0.74*10+1.3
			if res.Fare > (want + 0.1) ||  res.Fare < (want - 0.1) {
				t.Errorf("MakeEstimates(complex) returned fare= %f; want %f", res.Fare, want)
			}

		}else  if res.Id == 6 {
			want = 3.47
			if res.Fare > (want + 0.1) ||  res.Fare < (want - 0.1) {
				t.Errorf("MakeEstimates(complex) returned fare= %f; want %f", res.Fare, want)
			}

		}
	}

}	