package main

import (
	"encoding/json"
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

var mapValue = map[string]float32{"CDG": 20.3, "BOD": 24.7}

// fonction main
func main() {
	err := godotenv.Load("tempSensor.env")
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
	publish(client, "CDG")
	publish(client, "BOD")
	time.Sleep(10 * time.Second)
	repeatePublish(client)
}

func randomNumber(min float32, max float32) float32 {
	// Initialisation de la source de nombres aléatoires
	rand.Seed(time.Now().UnixNano())

	// Génération d'un nombre aléatoire compris entre min et max
	return min + rand.Float32()*(max-min)
}

func limitedRandomNumber(oldNumber float32) float32 {
	number := oldNumber + randomNumber(-0.2, 0.3)

	// Vérification des limites
	if number < -5 {
		number = -5
	} else if number > 44 {
		number = 44
	}

	return number
}

func saveValuesInJSON(value1 string, value2 string) []byte {
	// Création d'un objet map
	data := make(map[string]string)

	// Ajout des valeurs à l'objet map
	data["Date"] = value1
	data["Value"] = value2

	// Conversion de l'objet map en JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return jsonData
}

func publish(client mqtt.Client, airportcode string) {

	sensorType := 0
	sensorId := rand.Intn(2)

	floatAsString := strconv.FormatFloat(float64(mapValue[airportcode]), 'f', 2, 32)

	mapValue[airportcode] = limitedRandomNumber(mapValue[airportcode])

	fmt.Printf("Publishing message: %s to topic: %s\n", floatAsString, airportcode+"/sensors/"+strconv.Itoa(sensorType)+"/"+strconv.Itoa(sensorId))

	qos := byte(0)

	if os.Getenv("qos") == "1" {
		qos = byte(1)
	} else if os.Getenv("qos") == "2" {
		qos = byte(2)
	}

	var jsonData = saveValuesInJSON(time.Now().Format(time.RFC3339), floatAsString)

	token := client.Publish(airportcode+"/sensors/"+strconv.Itoa(sensorType)+"/"+strconv.Itoa(sensorId), qos, false, jsonData)
	token.Wait()
	time.Sleep(time.Second)
}
