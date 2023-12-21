package estimator

import (
	"beat.com/takehome/fareestimator/processor"
	"time"
	"sync"
)

type Estimate struct{
	Id		int64
	Fare	float64
}

func EstimateSegmentFare(input <-chan processor.SegmentData, output chan<- processor.SegmentData, wg *sync.WaitGroup)  {
	for inData := range input { 
		if inData.Speed <= 10 {
			inData.Fare = float64(11.90/(60*60)) * float64(inData.Epoch2 - inData.Epoch1)
		} else {
			if time.Unix(inData.Epoch1, 0).Hour() >= 0 && time.Unix(inData.Epoch1, 0).Hour() < 5 {
				if time.Unix(inData.Epoch2, 0).Hour() >= 0 && time.Unix(inData.Epoch2, 0).Hour() < 5 {
					inData.Fare = 1.30 * inData.Distance
				} else {
					//dT := inData.Epoch2 - inData.Epoch1
					//inData.Fare = ((float64(5-time.Unix(inData.Epoch1, 0).Hour())/float64(dT))* inData.Distance * 1.3) + ( ( float64(time.Unix(inData.Epoch2, 0).Hour()-5) / float64(dT) )* inData.Distance * 0.74 )
					inData.Fare = 1.30 * inData.Distance
				}
								
			} else {
				if time.Unix(inData.Epoch2, 0).Hour() >= 0 && time.Unix(inData.Epoch2, 0).Hour() < 5 {
					//dT := inData.Epoch2 - inData.Epoch1
					inData.Fare = 0.74 * inData.Distance
				} else {
					inData.Fare = 0.74 * inData.Distance
				}
				
			}
		}
		output<-inData
	}
	wg.Done()

}

func MakeEstimates (segsWithFare <-chan processor.SegmentData, out chan<- Estimate, wg *sync.WaitGroup){
	estimations := make(map[int64]float64)
	for result := range segsWithFare {
		if v, ok := estimations[result.Id];ok{
			v += result.Fare
			estimations[result.Id] = v

		} else {
			estimations[result.Id] = 1.3 + result.Fare	//First occurance of ID
		}
	}
	for k,v := range estimations{
		if v <  3.47 {
			estimations[k] = 3.47
		}

		out<-Estimate{Id: k, Fare: estimations[k] }
	}
	wg.Done()
}