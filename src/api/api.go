package main

import (
	db "aeroport/dbActions"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func initService() {
	db.InitDbClient()
}

func getBetweenDateTime(c *gin.Context) {
	sensortype := db.SensorType(0)

	if c.Query("sensor") == "1" {
		sensortype = db.TemperatureCel
	} else if c.Query("sensor") == "2" {
		sensortype = db.Atmospheric
	} else if c.Query("sensor") == "3" {
		sensortype = db.Pressure
	} else if c.Query("sensor") == "4" {
		sensortype = db.WindSpeed
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad sensor type"})
		return
	}

	d, err1 := time.Parse("2006-01-02T15:04", c.Query("from"))
	f, err2 := time.Parse("2006-01-02T15:04", c.Query("to"))
	if err1 != nil || err2 != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad date. " + err1.Error() + " | " + err2.Error()})
		return
	}

	res := db.GetMeasurementBetweenPeriod(
		sensortype,
		d,
		f,
	)

	c.IndentedJSON(http.StatusOK, res)
}

func getAverageForDay(c *gin.Context) {
	fmt.Println("GetAverageForDay")
	d, err := time.Parse("2006-01-02", c.Query("date"))
	if err != nil {
		fmt.Printf("value err : %s\n", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad date"})
		return
	}
	fmt.Println(d)

	res := db.GetAverageSensorsMeasurement(
		c.Query("airport"),
		d,
	)

	c.IndentedJSON(http.StatusOK, res)
}

func getDoc(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, []string{"Route 1 : GetBetweenDateTime?sensor=N&from=YYYY-MM-DDThh:mm:ss&to=YYYY-MM-DDThh:mm:ss", "Route 2 : GetAverageForDay?date=YYYY-MM-DD&airport=XXX"})
}

func main() {
	initService()

	router := gin.Default()

	router.GET("/GetBetweenDateTime", getBetweenDateTime)
	router.GET("/GetAverageForDay", getAverageForDay)
	router.GET("/", getDoc)

	router.Run("localhost:8080")
}
