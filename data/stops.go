package data

import (
	"github.com/gocarina/gocsv"
	"os"
)

type Stop struct {
	ID   string `csv:"stop_id"`
	Name string `csv:"stop_name"`
}

type StopReader struct {
	stops map[string]Stop
}

func NewStopReader() StopReader {
	var stops = []Stop{}
	var sr = StopReader{}

	stopsFile, err := os.OpenFile("./data/stops.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer stopsFile.Close()

	if err := gocsv.UnmarshalFile(stopsFile, &stops); err != nil {
		panic(err)
	}

	sr.stops = make(map[string]Stop, len(stops))

	for _, stop := range stops {
		sr.stops[stop.ID] = stop
	}

	return sr
}

func (sr StopReader) GetStop(id string) Stop {
	return sr.stops[id]
}

func (sr StopReader) Stops() []Stop {
	stops := make([]Stop, 0, len(sr.stops))
	for _, stop := range sr.stops {
		stops = append(stops, stop)
	}
	return stops
}
