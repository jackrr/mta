package api

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
)

type Stop struct {
	ID        string `csv:"stop_id"`
	Name      string `csv:"stop_name"`
	ParentID  string `csv:"parent_station"`
	stationID string
	ChildIDs  []string
}

type Station struct {
	ID        int
	StopIDs   []string
	Name      string
	transfers []Transfer
}

type Transfer struct {
	FromID string `csv:"from_stop_id"`
	ToID   string `csv:"to_stop_id"`
	Time   string `csv:"min_transfer_time"`
}

type StopReader struct {
	stops          map[string]Stop
	stations       map[int]Station
	stationsByName map[string]int
	nextStationID  int
}

func NewStopReader() StopReader {
	var sr = StopReader{}
	sr.loadStops()
	sr.initializeStations()
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

func (s Station) HasStop(stop Stop) bool {
	for _, stopID := range s.StopIDs {
		if stopID == stop.ID {
			return true
		}
	}

	return false
}

func (sr *StopReader) loadStops() {
	var stops = []Stop{}

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
		if stop.ParentID != "" {
			// Assumes parents always read in first
			parent := sr.stops[stop.ParentID]
			parent.ChildIDs = append(parent.ChildIDs, stop.ID)
			sr.stops[stop.ParentID] = parent
		}
	}
}

func (sr *StopReader) initializeStations() {
	var transfers = []Transfer{}

	transfersFile, err := os.OpenFile("./data/transfers.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer transfersFile.Close()

	if err := gocsv.UnmarshalFile(transfersFile, &transfers); err != nil {
		panic(err)
	}

	sr.stations = make(map[int]Station, len(sr.stops))
	sr.stationsByName = make(map[string]int, len(sr.stops))
	sr.nextStationID = 1
	for _, transfer := range transfers {
		fstop := sr.stops[transfer.FromID]
		tstop := sr.stops[transfer.ToID]
		if fstop.Name == tstop.Name {
			st := sr.getStationByName(fstop.Name)
			if st.Name == "" {
				st = sr.createStation(fstop.Name)
			}

			st.transfers = append(st.transfers, transfer)
			st.addStops([]Stop{fstop, tstop})
			st.addStops(sr.childStops(fstop))
			st.addStops(sr.childStops(tstop))
			sr.stations[st.ID] = st
		}
	}

	for _, station := range sr.stations {
		stopIDSets := findDisjointStopIDSets(station.transfers)
		if len(stopIDSets) > 1 {
			var stationToModify Station
			// split station into new station for each set
			for idx, stopIDList := range stopIDSets {
				if idx == 0 {
					// cut down existing station
					stationToModify = station
					stationToModify.StopIDs = []string{}
				} else {
					// create new station for any additional sets
					stationToModify = sr.createStation(fmt.Sprintf("%s (%s)", station.Name, string(stopIDList[0][0])))
				}

				for _, stopID := range stopIDList {
					stop := sr.GetStop(stopID)
					stationToModify.addStops([]Stop{stop})
					stationToModify.addStops(sr.childStops(stop))
				}
				sr.stations[stationToModify.ID] = stationToModify
			}
		}
	}
}

func (sr *StopReader) createStation(name string) (s Station) {
	s.ID = sr.nextStationID
	s.Name = name
	sr.stations[s.ID] = s
	sr.stationsByName[s.Name] = s.ID
	sr.nextStationID++
	return s
}

func (sr StopReader) getStationByName(name string) Station {
	return sr.stations[sr.stationsByName[name]]
}

func (s *Station) addStops(stops []Stop) {
	for _, stop := range stops {
		if s.HasStop(stop) {
			continue
		}
		s.StopIDs = append(s.StopIDs, stop.ID)
	}
}

func (sr StopReader) childStops(s Stop) (children []Stop) {
	for _, id := range s.ChildIDs {
		children = append(children, sr.GetStop(id))
	}

	return children
}

func (t Transfer) stopIDs() []string {
	return []string{t.FromID, t.ToID}
}

func findDisjointStopIDSets(transfers []Transfer) [][]string {
	sets := make([]map[string]bool, 0)

	for _, transfer := range transfers {
		// keep track of which sets we find the transfer in
		foundIn := make([]int, 0)

		for setIdx, set := range sets {
			if set[transfer.ToID] || set[transfer.FromID] {
				set[transfer.ToID] = true
				set[transfer.FromID] = true
				foundIn = append(foundIn, setIdx)
			}
		}

		if len(foundIn) > 1 {
			// merge found in sets down to one set
			newSets := make([]map[string]bool, 0)
			combinedSet := make(map[string]bool, 0)

			for setIdx, set := range sets {
				if contains(foundIn, setIdx) {
					for item := range set {
						combinedSet[item] = true
					}
				} else {
					newSets = append(newSets, set)
				}
			}

			sets = append(newSets, combinedSet)
		} else if len(foundIn) == 0 {
			// create new set and add transfer to it
			set := map[string]bool{transfer.ToID: true, transfer.FromID: true}
			sets = append(sets, set)
		}
	}

	results := [][]string{}
	for _, set := range sets {
		setArray := []string{}
		for item := range set {
			setArray = append(setArray, item)
		}
		results = append(results, setArray)
	}
	return results
}

func contains(nums []int, num int) bool {
	for _, n := range nums {
		if n == num {
			return true
		}
	}

	return false
}
