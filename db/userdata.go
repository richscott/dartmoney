package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for loading the Postgresql driver
)

// SecurityPosition holds the representation of the shares
// of a single particular security a client holds.
type SecurityPosition struct {
	Name   string `json:"name" db:"name"`
	Symbol string `json:"symbol" db:"symbol"`
	Shares int    `json:"shares" db:"shares"`
}

// UserPositions returns an array of the given user's securities
func UserPositions(emailAddr string) []SecurityPosition {
	db, err := sqlx.Connect("postgres", "user=finance dbname=portfolio sslmode=disable")
	if err != nil {
		log.Fatalln(err, "\nCould not connect to database - is Postgresql running?")
	}

	// Loop through rows using only one struct
	securityPos := SecurityPosition{}
	querySQL := `SELECT shares, symbol, name
							 FROM account_positions
							 WHERE account_id = (select id from account where email = $1)`

	positions := make([]SecurityPosition, 0)
	rows, err := db.Queryx(querySQL, emailAddr)
	if err != nil {
		log.Fatalln("error querying for user positions", err)
	}
	for rows.Next() {
		err := rows.StructScan(&securityPos)
		if err != nil {
			log.Fatalln(err)
		}
		positions = append(positions, securityPos)
	}

	fmt.Printf("positions = %#v\n", positions)
	return positions
}
