package main

import (
	db "aeroport/dbActions"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/measurement", measurementBetweenHandler)
	db.InitDbClient()

	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func postParams(r *http.Request, name string) string {
	err := r.ParseForm()
	if err != nil {
		log.Println("Bad form. " + err.Error())
	}
	return r.Form.Get(name)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("www/templates/index.gohtml")
	if err != nil {
		log.Println("Bad template. " + err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = tmpl.Execute(w, db.GetAllAirports())
}

func measurementBetweenHandler(w http.ResponseWriter, r *http.Request) {
	sensor, err := strconv.Atoi(postParams(r, "sensor"))
	if err != nil && sensor < 0 && sensor > 2 {
		log.Println("Bad sensor. " + err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sensorType := db.SensorType(sensor)

	airport := postParams(r, "airport")

	dateDebut, err1 := time.Parse("2006-01-02T15:04", postParams(r, "from"))
	dateFin, err2 := time.Parse("2006-01-02T15:04", postParams(r, "to"))
	if err1 != nil || err2 != nil {
		log.Println("Bad date. " + err1.Error() + " | " + err2.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	measures := db.GetMeasurementBetweenPeriod(sensorType, airport, dateDebut, dateFin)

	data := map[string]interface{}{
		"measures": measures,
	}

	switch sensor {
	case 0:
		data["Sensor"] = "Temperature"
	case 1:
		data["Sensor"] = "Pressure"
	case 2:
		data["Sensor"] = "WindSpeed"
	}

	tmpl, err := template.ParseFiles("www/templates/measurement.gohtml")
	if err != nil {
		log.Println("Bad template. " + err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Bad execute. " + err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
