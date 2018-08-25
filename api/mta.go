package api

import (
	"fmt"
	"github.com/jackrr/mta/data"
	"github.com/jackrr/mta/pb"
)

type MTA struct {
	f FeedGetter
}

func NewMTA(key string) MTA {
	mta := MTA{f: NewFeedGetter(key)}
	return mta
}

func (m MTA) GetTrain() {
	feed := m.f.GetFeed()
	sr := data.NewStopReader()
	r := data.NewRouteReader()
	var stop data.Stop
	var route data.Route
	var update *pb.TripUpdate
	var trip *pb.TripDescriptor

	for _, entity := range feed.GetEntity() {
		update = entity.GetTripUpdate()
		trip = update.GetTrip()
		route = r.GetRoute(trip.GetRouteId())

		for _, stu := range update.GetStopTimeUpdate() {
			stop = sr.GetStop(*stu.StopId)
			fmt.Printf("%v - %v\n", route.Name, stop.Name)
		}
	}
}
