package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/duckdb/duckdb-go/v2"
)

func main() {

	var geojson_sql = `
	WITH features AS (
        SELECT {
            'type': 'Feature',
            'geometry': ST_AsGeoJSON(ST_Point(longitudeDegrees, latitudeDegrees))::JSON,
            'properties': {
                'imo': imoNumber,
                'mmsi': mmsiNumber,
                'vesselName': vesselName,
                'sanction': CASE WHEN v.imo IS NOT NULL THEN true ELSE false END
            }
        } AS feature
        FROM mpa.ais_data m 
        LEFT JOIN vessel_sanction_imo_mmsi v ON m.imoNumber = v.imo 
        WHERE m.timeStamp::timestamp > (SELECT max(timeStamp::timestamp) - INTERVAL '5 minutes' FROM mpa.ais_data)
    )
    SELECT JSON_OBJECT(
        'type', 'FeatureCollection',
        'features', JSON_GROUP_ARRAY(feature)
    )::TEXT -- Cast to text for a raw dump
    FROM features;
	`

	db, err := sql.Open("duckdb", "duckdb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`ATTACH 'D:\workspace\mpa-golang\mpa-oceanx-web\database1.db' as mpa (TYPE SQLITE, READ_ONLY)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSTALL SPATIAL`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`LOAD SPATIAL`)
	if err != nil {
		log.Fatal(err)
	}

	/*
		_, err = db.Exec(`INSERT INTO people VALUES (42, 'John')`)
		if err != nil {
			log.Fatal(err)
		}

		var (
			id   int
			name string
		)
		row := db.QueryRow(`SELECT id, name FROM people`)
		err = row.Scan(&id, &name)
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("no rows")
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("id: %d, name: %s\n", id, name)
	*/

	var num int
	row := db.QueryRow(`SELECT  count(*) FROM mpa.ais_data;`)

	err = row.Scan(&num)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("no rows")
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("mpa.ais_data %d\n", num)

	row1 := db.QueryRow(`SELECT  count(*) FROM vessel_sanction_imo_mmsi;`)

	err = row1.Scan(&num)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("no rows")
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Vessel Sanction %d\n", num)
	var geojson_str string

	row2 := db.QueryRow(geojson_sql)
	err = row2.Scan(&geojson_str)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("no rows")
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Println(geojson_str)
}
