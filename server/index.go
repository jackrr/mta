package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackrr/mta/api"
	"net/http"
	"strconv"
)

func RunServer(mtaApiKey string) {
	m := api.NewMTA(mtaApiKey)
	router := gin.Default()

	router.GET("/stations", createSearchStationsHandler(&m))
	router.GET("/stations/:id/arrivals", createStationArrivalsHandler(&m))

	router.Run(":8001")
}

type stationsRequest struct {
	Query string `json:"query" form:"query" binding:"required"`
}

type stationsResponse struct {
	Stations []api.StationResponse `json:"stations"`
}

func createSearchStationsHandler(m *api.MTA) func(c *gin.Context) {
	return func(c *gin.Context) {
		var r stationsRequest
		c.Bind(&r)
		c.JSON(http.StatusOK, m.StationsMatching(r.Query))
	}
}

func createStationArrivalsHandler(m *api.MTA) func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, "Invalid station id specified")
		}
		c.JSON(http.StatusOK, m.UpcomingTrains(id))
	}
}
