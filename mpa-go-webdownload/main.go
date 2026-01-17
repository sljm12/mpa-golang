package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mpa"
	"net/http"
	"strings"

	"github.com/magiconair/properties"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//Load Properties
	p, err := properties.LoadFile("application.properties", properties.UTF8)
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	aisreports := get_mpa(p.GetString("apiKey", ""))
	print(aisreports)

	mpa.WriteAISReportToDB(aisreports, db)
}

func get_mpa(apiKey string) []mpa.AISReport {

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://sg-mdh-api.mpa.gov.sg/v1/vessel/positions/snapshot", nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("apikey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
	//reader := bufio.NewReader(resp.Body)

	byteValue, err := io.ReadAll(resp.Body)
	reports, err := unmarshal_stream(byteValue)

	fmt.Println(len(reports))

	return reports
}

func unmarshal_stream(byteValue []byte) ([]mpa.AISReport, error) {
	var reports []mpa.AISReport

	err := json.Unmarshal(byteValue, &reports)
	if err != nil {
		log.Fatalf("Error unmarshalling json: %v", err)
	}
	return reports, err
}

func print(reports []mpa.AISReport) {
	for _, s := range reports {
		cons := []string{s.Particulars.VesselName, s.TimeStamp}
		fmt.Println(strings.Join(cons, ","))
	}
}
