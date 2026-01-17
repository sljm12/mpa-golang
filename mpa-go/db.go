package mpa

import (
	"database/sql"
	"log"
)

// VesselParticulars represents the static details of the ship.
type VesselParticulars struct {
	VesselName    string  `json:"vesselName"`
	CallSign      string  `json:"callSign"`
	IMONumber     string  `json:"imoNumber"`
	Flag          string  `json:"flag"`
	VesselLength  float64 `json:"vesselLength"`
	VesselBreadth float64 `json:"vesselBreadth"`
	VesselDepth   float64 `json:"vesselDepth"`
	VesselType    string  `json:"vesselType"`
	GrossTonnage  float64 `json:"grossTonnage"`
	NetTonnage    float64 `json:"netTonnage"`
	Deadweight    float64 `json:"deadweight"`
	MMSINumber    string  `json:"mmsiNumber"`
	YearBuilt     string  `json:"yearBuilt"`
}

// AISReport represents a single vessel position report from the JSON array.
type AISReport struct {
	Particulars      VesselParticulars `json:"vesselParticulars"`
	Latitude         float64           `json:"latitude"`
	Longitude        float64           `json:"longitude"`
	LatitudeDegrees  float64           `json:"latitudeDegrees"`
	LongitudeDegrees float64           `json:"longitudeDegrees"`
	Speed            float64           `json:"speed"`
	Course           float64           `json:"course"`
	Heading          float64           `json:"heading"`
	DimA             int               `json:"dimA"`
	DimB             int               `json:"dimB"`
	TimeStamp        string            `json:"timeStamp"`
}

func WriteAISReportToDB(reports []AISReport, db *sql.DB) {
	// Define the SQL query matching the table structure we created earlier
	query := `
INSERT INTO ais_data (
    vesselName, callSign, imoNumber, flag, vesselLength, vesselBreadth, 
    vesselDepth, vesselType, grossTonnage, netTonnage, deadweight, 
    mmsiNumber, yearBuilt, latitude, longitude, latitudeDegrees, 
    longitudeDegrees, speed, course, heading, dimA, dimB, timeStamp
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Prepare the statement once
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Iterate over your data (assuming 'reports' is a []AISReport)
	for _, report := range reports {
		_, err := stmt.Exec(
			// Flattened VesselParticulars fields
			report.Particulars.VesselName,
			report.Particulars.CallSign,
			report.Particulars.IMONumber,
			report.Particulars.Flag,
			report.Particulars.VesselLength,
			report.Particulars.VesselBreadth,
			report.Particulars.VesselDepth,
			report.Particulars.VesselType,
			report.Particulars.GrossTonnage,
			report.Particulars.NetTonnage,
			report.Particulars.Deadweight,
			report.Particulars.MMSINumber,
			report.Particulars.YearBuilt,

			// Root level fields
			report.Latitude,
			report.Longitude,
			report.LatitudeDegrees,
			report.LongitudeDegrees,
			report.Speed,
			report.Course,
			report.Heading,
			report.DimA,
			report.DimB,
			report.TimeStamp,
		)

		if err != nil {
			log.Printf("Error inserting report for %s: %v", report.Particulars.VesselName, err)
		}
	}
}

func GetLatestReports(db *sql.DB) ([]AISReport, error) {
	query := `
        SELECT 
            vesselName, callSign, imoNumber, flag, vesselLength, vesselBreadth, 
            vesselDepth, vesselType, grossTonnage, netTonnage, deadweight, 
            mmsiNumber, yearBuilt, latitude, longitude, latitudeDegrees, 
            longitudeDegrees, speed, course, heading, dimA, dimB, timeStamp
        FROM ais_data 
        WHERE timeStamp >= datetime((SELECT MAX(timeStamp) FROM ais_data), '-2 minutes')`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []AISReport
	for rows.Next() {
		var r AISReport
		err := rows.Scan(
			&r.Particulars.VesselName,
			&r.Particulars.CallSign,
			&r.Particulars.IMONumber,
			&r.Particulars.Flag,
			&r.Particulars.VesselLength,
			&r.Particulars.VesselBreadth,
			&r.Particulars.VesselDepth,
			&r.Particulars.VesselType,
			&r.Particulars.GrossTonnage,
			&r.Particulars.NetTonnage,
			&r.Particulars.Deadweight,
			&r.Particulars.MMSINumber,
			&r.Particulars.YearBuilt,
			&r.Latitude,
			&r.Longitude,
			&r.LatitudeDegrees,
			&r.LongitudeDegrees,
			&r.Speed,
			&r.Course,
			&r.Heading,
			&r.DimA,
			&r.DimB,
			&r.TimeStamp,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}

	return reports, nil
}

func GetByIMO(db *sql.DB, imo string, datestart string, dateend string) ([]AISReport, error) {
	query := `
        SELECT 
            vesselName, callSign, imoNumber, flag, vesselLength, vesselBreadth, 
            vesselDepth, vesselType, grossTonnage, netTonnage, deadweight, 
            mmsiNumber, yearBuilt, latitude, longitude, latitudeDegrees, 
            longitudeDegrees, speed, course, heading, dimA, dimB, timeStamp
        FROM ais_data 
        WHERE imonumber = ? and timeStamp between datetime(?) and datetime (?))`

	rows, err := db.Query(query, imo, datestart, dateend)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return Rows2AISReport(rows)

}

func Rows2AISReport(rows *sql.Rows) ([]AISReport, error) {
	var reports []AISReport
	for rows.Next() {
		var r AISReport
		err := rows.Scan(
			&r.Particulars.VesselName,
			&r.Particulars.CallSign,
			&r.Particulars.IMONumber,
			&r.Particulars.Flag,
			&r.Particulars.VesselLength,
			&r.Particulars.VesselBreadth,
			&r.Particulars.VesselDepth,
			&r.Particulars.VesselType,
			&r.Particulars.GrossTonnage,
			&r.Particulars.NetTonnage,
			&r.Particulars.Deadweight,
			&r.Particulars.MMSINumber,
			&r.Particulars.YearBuilt,
			&r.Latitude,
			&r.Longitude,
			&r.LatitudeDegrees,
			&r.LongitudeDegrees,
			&r.Speed,
			&r.Course,
			&r.Heading,
			&r.DimA,
			&r.DimB,
			&r.TimeStamp,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}
