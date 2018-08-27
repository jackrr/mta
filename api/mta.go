package api

import (
	"fmt"
	"github.com/jackrr/mta/pb"
)

type MTA struct {
	f  FeedGetter
	sm StationManager
	r  RouteReader
}

type ExpectedArrival struct {
	Stop      string `json:"stop"`
	Train     string `json:"train"`
	Time      int64  `json:"time"`
	Direction string `json:"direction"`
}

type StationResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func NewMTA(key string) MTA {
	return MTA{f: NewFeedGetter(key), sm: NewStationManager(), r: NewRouteReader()}
}

func (m MTA) StationsMatching(query string) (res []StationResponse) {
	stations := m.sm.StationsMatching(query)
	for _, st := range stations {
		res = append(res, StationResponse{ID: st.ID, Name: st.Name})
	}
	return res
}

func (m MTA) UpcomingTrains(stationID int) []ExpectedArrival {
	var feed pb.FeedMessage
	var expectedArrivals []ExpectedArrival
	var route Route
	var update *pb.TripUpdate
	var trip *pb.TripDescriptor
	var arrivalTimeStamp int64
	var stop Stop

	station := m.sm.GetStation(stationID)
	fmt.Printf("%v", station)

	for _, feedID := range AllFeeds() {
		feed = m.f.GetFeed(feedID)
		for _, entity := range feed.GetEntity() {
			update = entity.GetTripUpdate()
			trip = update.GetTrip()
			route = m.r.GetRoute(trip.GetRouteId())

			for _, stu := range update.GetStopTimeUpdate() {
				stop = m.sm.GetStop(stu.GetStopId())

				if station.HasStop(stop) {
					arrivalTimeStamp = stu.GetArrival().GetTime()
					expectedArrivals = append(expectedArrivals, ExpectedArrival{
						Stop:      stop.Name,
						Train:     route.Name,
						Time:      arrivalTimeStamp,
						Direction: trip.GetDirection(),
					})
				}
			}
		}
	}

	return expectedArrivals
}
