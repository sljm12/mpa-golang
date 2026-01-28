package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/duckdb/duckdb-go/v2"
)

func main() {
	db, err := sql.Open("duckdb", "duckdb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`ATTACH 'D:\workspace\mpa-golang\mpa-oceanx-web\database1.db' as mpa (TYPE SQLITE, READ_ONLY)`)
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
	fmt.Printf("mpa.ais_data %d", num)

	row1 := db.QueryRow(`SELECT  count(*) FROM vessel_sanction_imo_mmsi;`)

	err = row1.Scan(&num)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("no rows")
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Vessel Sanction %d", num)
}
