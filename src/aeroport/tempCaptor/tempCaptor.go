package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

//A utiliser pour décoder le json par le sub aussi
//Contient les différents champs à récupérer
type captorDatas struct {
	IdCompteur   int     //id du capteur
	IdAeroport   string  //id aéroport (code "IATA" 3 caractères)
	NatureMesure string  //Nature mesure ("Temperature","Atmospheric pressure", "Wind speed")
	Valeur       float32 //Valeur de la mesure (numérique)
	DateHeureMes string  //Date et heure de la mesure (timestamp : YYYY-MM-DD-hh-mm-ss)
}

//Créer client
func createClientOptions(brokerUrl string, clientId string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerUrl)
	opts.SetClientID(clientId)
	return opts
}

//Connection au broker
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

//fonction random
func generateDatas(codeAeroport string) captorDatas {
	t := time.Now()
	var data = captorDatas{
		IdCompteur:   1,
		IdAeroport:   codeAeroport,
		NatureMesure: "Temperature",
		DateHeureMes: t.Format("2006-01-02-15-04-05"),
	}
	data.Valeur = randomTemps(data.IdAeroport)
	return data
}

//TODO
func randomTemps(codeAer string) float32 {
	return 1.1
}

//Fonction qui génère une string au format JSON à partir d'une structure captorDatas
func encodeJson(datas captorDatas) string {
	empData := &datas
	e, err := json.Marshal(empData)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(e)
}

//fonction main
func main() {

	//var codesAeroport = []string{"CDG", "BOD", "CFE", "DIJ", "GNB", "JCA", "LRH", "NTE"}

	//A modifier --> fonctionne en local
	urlBroker := "tcp://localhost:1883"
	clientId2 := "clientVic_pub"

	client_pub := connect(urlBroker, clientId2)

	for {
		data := generateDatas("CDG") //CDG à modif dans config réelle
		jsonDatas := encodeJson(data)

		client_pub.Publish("test", 0, false, jsonDatas) //topic à changer
		time.Sleep(10 * time.Second)
		if !client_pub.IsConnected() {
			break
		}
	}

}
