package logActions

import (
	db "aeroport/dbActions"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const logPath = "./Logs"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkIfFolderExist(folder string) bool {
	if _, err := os.Stat(folder); !os.IsNotExist(err) {
		return true
	}
	return false
}

func getExistingFileDatas(fileUri string) []byte {
	dat, err := os.ReadFile(fileUri)

	if err != nil {
		file, err := os.Create(fileUri)
		if err != nil {
			fmt.Println("Err during creating log file")
			log.Fatal(err)
		}
		file.Close()
		return []byte("Airport, Date, Value, Captor, SensorType\n")
	}
	return dat
}

func formatSensMeasurementToCsv(data db.SensorMeasurement) string {
	csv := ""
	csv += data.Airport + ","
	csv += data.Datetime.Time().Format("2006-01-02 15:04:05") + ","
	csv += strconv.FormatFloat(data.Value, 'f', 2, 64) + ","
	csv += strconv.Itoa(data.Captor) + ","
	csv += strconv.Itoa(int(data.Sensortype)) + "\n"
	return csv
}

func writeDatas(data db.SensorMeasurement) bool {
	date := time.Now().Format("2006-01-02")

	err := os.WriteFile(logPath+"/"+data.Airport+"/"+date+".csv",
		append(getExistingFileDatas(logPath+"/"+data.Airport+"/"+date+".csv"), []byte(formatSensMeasurementToCsv(data))[:]...), 0644)

	if err != nil {
		fmt.Printf("Error while writing csv file: %v", err)
		return false
	}
	return true

}

func WriteLog(datas db.SensorMeasurement) bool {
	if !checkIfFolderExist(logPath) {
		err := os.Mkdir(logPath, 0777)
		if err != nil {
			fmt.Println("Error while creating log folder")
			log.Fatal(err)
		}
	}

	if !checkIfFolderExist(logPath + "/" + datas.Airport) {
		err := os.Mkdir(logPath+"/"+datas.Airport, 0777)
		if err != nil {
			fmt.Printf("Error while creating log day folder: %v", err)
			return false
		}
	}

	return writeDatas(datas)
}
