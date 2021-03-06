package api

import (
	"bytes"
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"regexp"
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

type StationManager struct {
	stops          map[string]Stop
	stations       map[int]Station
	stationsByName map[string]int
	nextStationID  int
}

func NewStationManager() StationManager {
	var sr = StationManager{}
	sr.loadStops()
	sr.initializeStations()
	return sr
}

func (sr StationManager) GetStop(id string) Stop {
	return sr.stops[id]
}

func (sr StationManager) GetStation(id int) Station {
	return sr.stations[id]
}

func (sr StationManager) Stops() []Stop {
	stops := make([]Stop, 0, len(sr.stops))
	for _, stop := range sr.stops {
		stops = append(stops, stop)
	}
	return stops
}

func (sr StationManager) StationsMatching(query string) (stations []Station) {
	re, err := regexp.Compile(fmt.Sprintf("(?i).*%s.*", query))
	if err != nil {
		return stations
	}

	for name, id := range sr.stationsByName {
		if re.MatchString(name) {
			stations = append(stations, sr.stations[id])
		}
	}
	return stations
}

func (s Station) HasStop(stop Stop) bool {
	for _, stopID := range s.StopIDs {
		if stopID == stop.ID {
			return true
		}
	}

	return false
}

func (sr *StationManager) loadStops() {
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

func (sr *StationManager) initializeStations() {
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
			sr.registerStation(st)
		}
	}

	updatedStations := []Station{}
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
					stationToModify = sr.createStation(station.Name)
				}

				for _, stopID := range stopIDList {
					stop := sr.GetStop(stopID)
					stationToModify.addStops([]Stop{stop})
					stationToModify.addStops(sr.childStops(stop))
				}

				stationToModify.computeName()
				updatedStations = append(updatedStations, stationToModify)
			}
		} else {
			station.computeName()
			updatedStations = append(updatedStations, station)
		}
	}

	sr.resetStations()

	for _, updatedStation := range updatedStations {
		sr.registerStation(updatedStation)
	}
}

func (sr *StationManager) createStation(name string) (s Station) {
	s.ID = sr.nextStationID
	s.Name = name
	sr.nextStationID++
	return s
}

func (sr *StationManager) registerStation(s Station) {
	sr.stations[s.ID] = s
	sr.stationsByName[s.Name] = s.ID
}

func (sr *StationManager) resetStations() {
	sr.stations = map[int]Station{}
	sr.stationsByName = map[string]int{}
}

func (sr StationManager) getStationByName(name string) Station {
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

func (s *Station) computeName() {
	lines := map[rune]bool{}
	for _, stopID := range s.StopIDs {
		lines[rune(stopID[0])] = true
	}

	var name bytes.Buffer

	name.WriteString(s.Name)
	name.WriteString(" (")
	idx := 0
	for line := range lines {
		name.WriteRune(line)
		if idx < len(lines)-1 {
			name.WriteRune(',')
		}
		idx++
	}
	name.WriteRune(')')
	s.Name = name.String()
}

func (sr StationManager) childStops(s Stop) (children []Stop) {
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
