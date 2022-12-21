package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Context init
var collection *mongo.Collection
var ctx = context.TODO()

// Connection to the database
func initDbClient() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://mqttAirportSub:mqttAirportSub99@cluster0.vp9lmsa.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("airport").Collection("sensormeasurements")
}

type SensorType int

const (
	TemperatureCel = 1
	Atmospheric    = 2
	Pressure       = 3
	WindSpeed      = 4
)

type SensorMeasurement struct {
	Captor     int
	Airport    string
	Sensortype SensorType
	Value      float32
	Datetime   primitive.DateTime
}

type SensorAverageMeasurement struct {
	Sensortype SensorType
	Value      float32
	Count      int
}

// Insert a new measurement in the database
func AddValue(measure SensorMeasurement) {
	fmt.Printf("Inserting: %v\n", measure)
	result, insertError := collection.InsertOne(ctx, measure)
	if insertError != nil {
		log.Fatal(insertError)
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}

func GetMeasurementBetweenPeriod(sensorType SensorType, start time.Time, end time.Time) []SensorMeasurement {
	fmt.Printf("Start: %v\n", primitive.NewDateTimeFromTime(start))
	filter := bson.M{
		"sensortype": sensorType,
		"datetime": bson.M{
			"$gte": primitive.NewDateTimeFromTime(start),
			"$lte": primitive.NewDateTimeFromTime(end),
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	var results []SensorMeasurement
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	return results
}

func GetAverageSensorsMeasurement(airport string, date time.Time) []SensorAverageMeasurement {
	y, m, d := date.Date()
	startDate := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	endDate := time.Date(y, m, d, 23, 59, 59, 999, time.Local)

	filter := bson.M{
		"airport": airport,
		"datetime": bson.M{
			"$gte": primitive.NewDateTimeFromTime(startDate),
			"$lte": primitive.NewDateTimeFromTime(endDate),
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	var results []SensorMeasurement
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	var averages = []SensorAverageMeasurement{
		SensorAverageMeasurement{
			Sensortype: TemperatureCel,
			Value:      0,
			Count:      0,
		},
		SensorAverageMeasurement{
			Sensortype: Atmospheric,
			Value:      0,
			Count:      0,
		},
		SensorAverageMeasurement{
			Sensortype: Pressure,
			Value:      0,
			Count:      0,
		},
		SensorAverageMeasurement{
			Sensortype: WindSpeed,
			Value:      0,
			Count:      0,
		},
	}

	for _, result := range results {
		averages[result.Sensortype].Value += result.Value
		averages[result.Sensortype].Count++
	}

	for i := range averages {
		averages[i].Value = averages[i].Value / float32(averages[i].Count)
	}

	return averages
}
