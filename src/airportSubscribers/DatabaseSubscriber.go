package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var collection *mongo.Collection
var ctx = context.TODO()

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://mqttAirportSub:mqttAirportSub99@cluster0.vp9lmsa.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Par aérport : Valeur Vente, Temp et pression, date, Id Capteur
