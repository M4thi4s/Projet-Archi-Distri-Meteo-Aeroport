package main

import (
	db "aeroport/dbActions"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func initService() {
	db.InitDbClient()
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	keepAlive := make(chan os.Signal)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)

	fmt.Printf("Starting service\n")

	initService()
	/*

		var startTime = time.Now().Add(-time.Minute * 50)
		var endTime = time.Now().Add(time.Minute * 10)

		var res1 = GetMeasurementBetweenPeriod(Pressure, startTime, endTime)
		fmt.Printf("Result 1: %v\n", res1)

		var res2 = GetAverageSensorsMeasurement("AAA", time.Now())
		fmt.Printf("Result 2: %v\n", res2)
	*/

	var broker = "51.210.45.234"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("DatabaseClient")
	opts.SetUsername("DatabaseClient")
	opts.SetPassword("DbPass123")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)
	publish(client)

	client.Disconnect(250)
	<-keepAlive
}

func publish(client mqtt.Client) {
	text := fmt.Sprintf("9.85")
	token := client.Publish("AIR/sensors/1/1234", 0, false, text)
	token.Wait()
	time.Sleep(time.Second)
}

func onMessage(c mqtt.Client, msg mqtt.Message) {
	topicDatas := strings.Split(msg.Topic(), "/")

	sensortype := db.SensorType(0)

	if topicDatas[2] == "1" {
		sensortype = db.TemperatureCel
	} else if topicDatas[2] == "2" {
		sensortype = db.Atmospheric
	} else if topicDatas[2] == "3" {
		sensortype = db.Pressure
	} else if topicDatas[2] == "4" {
		sensortype = db.WindSpeed
	} else {
		fmt.Printf("Bad sensor type %s", topicDatas[1])
		return
	}

	captorId, errCaptorId := strconv.Atoi(topicDatas[3])

	if errCaptorId != nil {
		fmt.Printf("Captor Id err : %s", errCaptorId)
		return
	}

	sensorVal, errValue := strconv.ParseFloat(string(msg.Payload()), 64)

	if errValue != nil {
		fmt.Printf("value err : %s", errValue)
		return
	}

	newMeasure := db.SensorMeasurement{
		Captor:     captorId,
		Airport:    topicDatas[0],
		Sensortype: sensortype,
		Value:      sensorVal,
		Datetime:   primitive.NewDateTimeFromTime(time.Now()),
	}
	db.AddValue(newMeasure)
}

func sub(client mqtt.Client) {
	topic := "+/sensors/+/+"
	token := client.Subscribe(topic, 1, onMessage)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
