package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

// fonction main
func main() {
	err := godotenv.Load("pressureSensor.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var broker = os.Getenv("broker")
	port, portErr := strconv.Atoi(os.Getenv("port"))

	if portErr != nil {
		log.Fatalf("Error loading env port value")
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(os.Getenv("id"))
	opts.SetUsername(os.Getenv("username"))
	opts.SetPassword(os.Getenv("psw"))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	repeatePublish(client)

	client.Disconnect(250)
}

func repeatePublish(client mqtt.Client) {
	publish(client)
	time.Sleep(10 * time.Second)
	repeatePublish(client)
}

func publish(client mqtt.Client) {
	var codesAeroport = []string{"CDG", "BOD"}

	randomCode := codesAeroport[rand.Intn(len(codesAeroport))]

	sensorType := 1
	sensorId := rand.Intn(2)

	float := rand.Float32()*100 - 50

	floatAsString := strconv.FormatFloat(float64(float), 'f', 2, 32)

	fmt.Printf("Publishing message: %s to topic: %s\n", floatAsString, randomCode+"/sensors/"+strconv.Itoa(sensorType)+"/"+strconv.Itoa(sensorId))

	qos := byte(0)

	if os.Getenv("qos") == "1" {
		qos = byte(1)
	} else if os.Getenv("qos") == "2" {
		qos = byte(2)
	}

	token := client.Publish(randomCode+"/sensors/"+strconv.Itoa(sensorType)+"/"+strconv.Itoa(sensorId), qos, false, floatAsString)
	token.Wait()
	time.Sleep(time.Second)
}
