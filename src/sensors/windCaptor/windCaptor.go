package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/rand"
	"time"
)

// A utiliser pour décoder le json par le sub aussi
// Contient les différents champs à récupérer
type captorDatas struct {
	IdSensor     int     //id du capteur
	IdAeroport   string  //id aéroport (code "IATA" 3 caractères)
	NatureMesure string  //Nature mesure ("Temperature","Atmospheric pressure", "Wind speed")
	Valeur       float32 //Valeur de la mesure (numérique)
	DateHeureMes string  //Date et heure de la mesure (timestamp : YYYY-MM-DD-hh-mm-ss)
}

// Créer client
func createClientOptions(brokerUrl string, clientId string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerUrl)
	opts.SetClientID(clientId)
	return opts
}

// Connection au broker
func connect(brokerURI string, clientId string) mqtt.Client {
	fmt.Println("Trying to connect (" + brokerURI + ", " + clientId + ")...")
	opts := createClientOptions(brokerURI, clientId)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

// fonction random
func generateDatas(codeAeroport string) captorDatas {
	t := time.Now()
	var data = captorDatas{
		IdSensor:     2,
		IdAeroport:   codeAeroport,
		NatureMesure: "Wind speed",
		DateHeureMes: t.Format("2006-01-02-15-04-05"),
	}
	data.Valeur = randomWind()
	return data
}

// TODO: à améliorer pour prendre en compte les anciennes valeurs de chaque airport pour pas trop faire varier les résultats
func randomWind() float32 {
	return float32(rand.Intn(80))
}

// Fonction qui génère une string au format JSON à partir d'une structure captorDatas
func encodeJson(datas captorDatas) string {
	empData := &datas
	e, err := json.Marshal(empData)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(e)
}

// fonction main
func main() {

	var codesAeroport = []string{"CDG", "BOD", "CFE", "DIJ", "GNB", "JCA", "LRH", "NTE"}

	//TODO: à changer lors du fichier config
	urlBroker := "tcp://51.210.45.234:1883"
	clientId := "windCaptor"

	client_pub := connect(urlBroker, clientId)

	for _, airportCode := range codesAeroport {
		data := generateDatas(airportCode)
		jsonDatas := encodeJson(data)

		client_pub.Publish(urlBroker+"/airport/"+"airportCode"+"/sensor/wind", 0, false, jsonDatas) //topic à changer
		time.Sleep(3 * time.Minute)
		if !client_pub.IsConnected() {
			break
		}
	}

}
