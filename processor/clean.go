package processor

import "beat.com/takehome/fareestimator/inputfeeder"
import "sync"

func PrepareAndClean(input <-chan inputfeeder.Datarow, output chan<- SegmentData){
	dataSeqPerId := make(map[int64] chan inputfeeder.Datarow)
	var wgClean sync.WaitGroup
	
	for inData := range input {
		if c, exist := dataSeqPerId[inData.Id]; exist{
			c <- inData
		} else {
			dataSeqPerId[inData.Id] = make(chan inputfeeder.Datarow, 100000)
			dataSeqPerId[inData.Id] <- inData
			wgClean.Add(1)
			go Clean(dataSeqPerId[inData.Id], output, &wgClean)
		}
	}
	for _,c := range dataSeqPerId{
		close(c)
	}
	defer close(output)
	wgClean.Wait()
}

func Clean(input <-chan inputfeeder.Datarow, output chan<- SegmentData, wg *sync.WaitGroup){
	segments := make(map[int64]SegmentData)
	for inData := range input {
		if s, ok := segments[inData.Id];ok{
			s.Epoch1 = s.Epoch2
			s.Lat1 = s.Lat2
			s.Lon1 = s.Lon2
			s.Epoch2 = inData.Epoch
			s.Lat2 = inData.Lat
			s.Lon2 = inData.Lon
			s.Distance, s.Speed = CalcDistanceSpeed(s)
			if s.Speed > 100 {continue}
			segments[inData.Id] = s
			output <-s
		}else{
			segments[inData.Id] = SegmentData{
				Id:       inData.Id,
				Epoch2:   inData.Epoch,
				Lat2:     inData.Lat,
				Lon2:     inData.Lon,
			}

		}

	}
	wg.Done()

}