package inputfeeder

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)
type Datarow struct{
	Id int64
	Lat, Lon float64
	Epoch int64
}

func StreamReader(inputFilePath string, out chan<- Datarow) {
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalln("Failed to open input file. Closing the program...", err)
	}
	defer f.Close()
	log.Println("Input file successfuly opened: ", inputFilePath)
	csvReader := csv.NewReader(f)
	start := time.Now()
	/*recs, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Reading csv failed" ,err)
	}*/
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			log.Println("Reached the end of input file. Stopping...")
			break
		}
		if err != nil {
			log.Fatalln("Could not open input file Closing the program... ", err)
		}

		id, err := strconv.ParseInt(rec[0], 10, 64)
		if err != nil {
			log.Println("Could not read id column. Ingonroing input row... ", err)
			continue
		}
		lat, err := strconv.ParseFloat(rec[1], 64)
		if err != nil {
			log.Println("Could not read lat column. Ingonroing input row... ", err)
			continue
		}
		lon, err := strconv.ParseFloat(rec[2], 64)
		if err != nil {
			log.Println("Could not read lon column. Ingonroing input row... ", err)
			continue
		}
		epoch, err := strconv.ParseInt(rec[3], 0, 64)
		if err != nil {
			log.Println("Could not read timestamp column. Ingonroing input row... ", err)
			continue
		}
		out <- Datarow{Id: id, Lat: lat, Lon: lon, Epoch: epoch}
	}
	log.Println("Time reading input data: ", time.Since(start))
	close(out) //No more data, signaling

}