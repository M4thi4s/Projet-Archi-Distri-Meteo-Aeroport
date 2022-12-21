package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func initService() {
	initDbClient()
}

func main() {
	fmt.Printf("Starting service\n")

	initService()

	sendtest := SensorMeasurement{
		Captor:     2,
		Airport:    "AAA",
		Sensortype: Pressure,
		Value:      1013.0,
		Datetime:   primitive.NewDateTimeFromTime(time.Now()),
	}
	AddValue(sendtest)

	var startTime = time.Now().Add(-time.Minute * 50)
	var endTime = time.Now().Add(time.Minute * 10)

	var res1 = GetMeasurementBetweenPeriod(Pressure, startTime, endTime)
	fmt.Printf("Result 1: %v\n", res1)

	var res2 = GetAverageSensorsMeasurement("AAA", time.Now())
	fmt.Printf("Result 2: %v\n", res2)
}
