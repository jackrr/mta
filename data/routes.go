package data

import (
	"github.com/gocarina/gocsv"
	"os"
)

type Route struct {
	ID       string `csv:"route_id"`
	Name     string `csv:"route_short_name"`
	LongName string `csv:"route_long_name"`
}

type RouteReader struct {
	routes map[string]Route
}

func NewRouteReader() RouteReader {
	var routes = []Route{}
	var r = RouteReader{}

	routesFile, err := os.OpenFile("./data/routes.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer routesFile.Close()

	if err := gocsv.UnmarshalFile(routesFile, &routes); err != nil {
		panic(err)
	}

	r.routes = make(map[string]Route, len(routes))

	for _, route := range routes {
		r.routes[route.ID] = route
	}

	return r
}

func (r RouteReader) GetRoute(id string) Route {
	return r.routes[id]
}

func (r RouteReader) Routes() []Route {
	routes := make([]Route, 0, len(r.routes))
	for _, route := range r.routes {
		routes = append(routes, route)
	}
	return routes
}
