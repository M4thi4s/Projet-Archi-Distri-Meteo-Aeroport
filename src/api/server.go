package main

import (
	db "aeroport/dbActions"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func initService() {
	db.InitDbClient()
}

func getBetweenDateTime(c *gin.Context) {
	sensortype := db.SensorType(0)

	if c.Param("sensor") == "0" {
		sensortype = db.TemperatureCel
	} else if c.Param("sensor") == "1" {
		sensortype = db.Pressure
	} else if c.Param("sensor") == "2" {
		sensortype = db.WindSpeed
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad sensor type"})
		return
	}

	//

	d, err1 := time.Parse("2006-01-02T15:04", c.Query("from"))
	f, err2 := time.Parse("2006-01-02T15:04", c.Query("to"))
	if err1 != nil || err2 != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad date. " + err1.Error() + " | " + err2.Error()})
		return
	}

	res := db.GetMeasurementBetweenPeriod(
		sensortype,
		c.Param("airport"),
		d,
		f,
	)

	c.IndentedJSON(http.StatusOK, res)
}

func getAverageForDay(c *gin.Context) {
	fmt.Println("GetAverageForDay")
	d, err := time.Parse("2006-01-02", c.Param("date"))
	if err != nil {
		fmt.Printf("value err : %s\n", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad date"})
		return
	}
	fmt.Println(d)

	res := db.GetAverageSensorsMeasurement(
		c.Param("airport"),
		d,
	)

	c.IndentedJSON(http.StatusOK, res)
}

func getDoc(c *gin.Context) {
	jsonFile, err := os.ReadFile("api/openapi.yaml")
	if err != nil {
		fmt.Println("File reading error", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Unable to find open API file"})
		return
	} else {
		c.Header("Content-Disposition", "attachment; filename=openAPI.yaml")
		c.Data(http.StatusOK, "application/octet-stream", jsonFile)
	}
}

func main() {
	initService()

	router := gin.Default()

	router.GET("/GetBetweenDateTime/:airport/:sensor", getBetweenDateTime)
	router.GET("/GetAverageForDay/:airport/:date", getAverageForDay)
	router.GET("/", getDoc)

	router.Run("localhost:8080")
}
