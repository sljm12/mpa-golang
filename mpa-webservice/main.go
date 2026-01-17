package main

import (
	"database/sql"
	"fmt"
	"log"
	"mpa"
	"net/http"

	"github.com/magiconair/properties"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Handlers struct {
	DB *sql.DB
}

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reports, err := mpa.GetLatestReports(h.DB)

	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Error in request")
	}

	featureCollection := AISReportToGeoJson(reports)
	byteValue, err := featureCollection.MarshalJSON()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, string(byteValue))
}

type GetByIMO struct {
	Handlers
}

/*
	func (h *GetByIMO) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		imo := r.URL.Query().Get("imo")
		ds := r.URL.Query().Get("datestart")
		de := r.URL.Query().Get("dateend")
		reports, err = mpa.GetByIMO(h.DB, imo, ds, de)

		if err != nil {
			log.Fatal(err)
			fmt.Fprintf(w, "Error in request")
		}

		featureCollection := AISReportToGeoJson(reports)
		byteValue, err := featureCollection.MarshalJSON()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, string(byteValue))
	}
*/
func main() {
	p, err := properties.LoadFile("application.properties", properties.UTF8)
	if err != nil {
		fmt.Println(err)
	}
	db, err := sql.Open("sqlite3", p.GetString("db", ""))

	if err != nil {
		log.Fatal(err)
	}

	handlers := &Handlers{
		DB: db,
	}

	//reports, err := mpa.GetLatestReports(db)

	//j_str, err := json.Marshal(&reports)

	//fmt.Println(string(j_str))
	/*
		for _, s := range reports {
			fmt.Println(s.Particulars.VesselName)
		}
	*/
	/*
		featureCollection := AISReportToGeoJson(reports)
		byteValue, err := featureCollection.MarshalJSON()
		fmt.Println(string(byteValue))
	*/

	http.Handle("/GetLatest", handlers)
	fmt.Println("Server listening on :8000")
	http.ListenAndServe(":8000", nil)
}

func LatestHandler(w http.ResponseWriter, r *http.Request) {

}

func AISReportToGeoJson(reports []mpa.AISReport) geojson.FeatureCollection {

	featureCollection := geojson.NewFeatureCollection()

	for _, r := range reports {
		feature := geojson.NewFeature(orb.Point{r.LongitudeDegrees, r.LatitudeDegrees})
		feature.Properties["vesselName"] = r.Particulars.VesselName
		feature.Properties["Speed"] = r.Speed
		feature.Properties["Course"] = r.Course
		feature.Properties["VesselType"] = r.Particulars.VesselType
		feature.Properties["IMO"] = r.Particulars.IMONumber
		feature.Properties["MMSI"] = r.Particulars.MMSINumber
		feature.Properties["TimeSTamp"] = r.TimeStamp
		featureCollection.Append(feature)
	}
	return *featureCollection
}
