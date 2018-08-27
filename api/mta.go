package api

import (
	"fmt"
	"github.com/jackrr/mta/pb"
	"time"
)

type MTA struct {
	f  FeedGetter
	sm StationManager
}

type ExpectedArrival struct {
	Stop      Stop
	Route     Route
	Time      time.Time
	Direction string
}

type StationResponse struct {
	Name string
	ID   int
}

func NewMTA(key string) MTA {
	mta := MTA{f: NewFeedGetter(key), sm: NewStationManager()}
	return mta
}

func (m MTA) StationsMatching(query string) (res []StationResponse) {
	stations := m.sm.StationsMatching(query)
	for _, st := range stations {
		res = append(res, StationResponse{ID: st.ID, Name: st.Name})
	}
	return res
}

func (m MTA) UpcomingTrains(stationName string) []string {
	arrivals := m.expectedArrivals(stationName)
	updates := make([]string, len(arrivals))
	for i, arrival := range arrivals {
		updates[i] = arrival.String()
	}

	return updates
}

func (a ExpectedArrival) String() string {
	return fmt.Sprintf("%v - %v (%v -- %v)\n", a.Route.Name, a.Time, a.Stop.Name, a.Stop.ID)
}

func (m MTA) expectedArrivals(stationName string) []ExpectedArrival {
	var feed pb.FeedMessage
	var expectedArrivals []ExpectedArrival
	var route Route
	var update *pb.TripUpdate
	var trip *pb.TripDescriptor
	var arrivalTimeStamp int64
	var stop Stop

	r := NewRouteReader()
	station := m.sm.getStationByName(stationName)
	fmt.Printf("%v", station)

	for _, feedID := range AllFeeds() {
		feed = m.f.GetFeed(feedID)
		for _, entity := range feed.GetEntity() {
			update = entity.GetTripUpdate()
			trip = update.GetTrip()
			route = r.GetRoute(trip.GetRouteId())

			for _, stu := range update.GetStopTimeUpdate() {
				stop = m.sm.GetStop(stu.GetStopId())

				if station.HasStop(stop) {
					arrivalTimeStamp = stu.GetArrival().GetTime()
					expectedArrivals = append(expectedArrivals, ExpectedArrival{
						Stop:      stop,
						Route:     route,
						Time:      time.Unix(arrivalTimeStamp, 0),
						Direction: trip.GetDirection(),
					})
				}
			}
		}
	}

	return expectedArrivals
}
