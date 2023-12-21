package main

import (
	"fmt"
	"sync"
	"log"
	"os"
	"time"
	"runtime"
	"strconv"
	"encoding/csv"
	"beat.com/takehome/fareestimator/estimator"
	"beat.com/takehome/fareestimator/inputfeeder"
	"beat.com/takehome/fareestimator/processor"
)

const BUFFER = 100000
func main(){
	start := time.Now()
	cpuN := runtime.NumCPU()
	inRows := make(chan inputfeeder.Datarow, BUFFER)
	go inputfeeder.StreamReader(os.Args[1], inRows)

	segs :=make(chan processor.SegmentData, BUFFER)
	go processor.PrepareAndClean(inRows, segs)

	var wgFare sync.WaitGroup
	wgFare.Add(cpuN)
	segsWithFare :=make(chan processor.SegmentData, BUFFER)
	for i:=1; i <= cpuN; i++ {go estimator.EstimateSegmentFare(segs, segsWithFare, &wgFare)}
	go func() {
        defer close(segsWithFare)
        wgFare.Wait()
    }()
	
	var wgTotal sync.WaitGroup
	wgTotal.Add(1)
	finalEstimates := make(chan estimator.Estimate, BUFFER)
	go estimator.MakeEstimates(segsWithFare, finalEstimates, &wgTotal)	
	go func() {
        defer close(finalEstimates)
        wgTotal.Wait()
    }()
		
	csvFile, err := os.Create(os.Args[2])

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	} else {
		defer csvFile.Close()
	}
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Write([]string{"id_ride", "fare_estimate" })

	for est := range finalEstimates {
		fmt.Printf("id_ride %d fare_estimate %9.2f\n", est.Id, est.Fare)
		csvwriter.Write([]string{strconv.FormatInt(est.Id,10), strconv.FormatFloat(est.Fare,'f',2,64) })
	} 
	log.Println("Total Time: ", time.Since(start))


	csvwriter.Flush()
	
}


