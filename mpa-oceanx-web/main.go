package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mpa"
	"net/http"

	"github.com/magiconair/properties"
	_ "github.com/mattn/go-sqlite3"
)

func unmarshal_stream(byteValue []byte) ([]mpa.AISReport, error) {
	var reports []mpa.AISReport

	err := json.Unmarshal(byteValue, &reports)
	if err != nil {
		log.Fatalf("Error unmarshalling json: %v", err)
	}
	return reports, err
}

func get_mpa(apiKey string) []mpa.AISReport {

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://oceans-x.mpa.gov.sg/api/v1/vessel/positions/1.0.0/snapshot", nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("ApiKey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
	//reader := bufio.NewReader(resp.Body)

	byteValue, err := io.ReadAll(resp.Body)
	aisReports, err := unmarshal_stream(byteValue)

	return aisReports
}

func main() {
	p, err := properties.LoadFile("application.properties", properties.UTF8)
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("sqlite3", "./database1.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println(p.GetString("apiKey", ""))

	data := get_mpa(p.GetString("apiKey", ""))
	mpa.WriteAISReportToDB(data, db)

}
