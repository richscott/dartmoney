package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const schema = `
DROP TABLE IF EXISTS account_positions;
DROP TABLE IF EXISTS account;

CREATE TABLE account (
    id serial primary key,
    email varchar(120) unique,
    hash_pass text not null,
    name varchar(120)
);

CREATE TABLE account_positions (
    symbol varchar(24) not null,
    name varchar(120) not null,
    shares int not null,
    account_id int references account(id)
)`

// CreateSchema creates a new finance portfolio schema
func CreateSchema() {
	// this Pings the database trying to connect, panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("postgres", "user=finance dbname=portfolio sslmode=disable")
	if err != nil {
		log.Fatalln(err, "\nCould not connect to database - is Postgresql running?")
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO account (email, hash_pass, name) VALUES ($1, $2, $3)",
		"investor@somewhere", "xxx", "Hopeful Investor")
	tx.Commit()

	tx = db.MustBegin()
	sharesInsert := `INSERT INTO account_positions (symbol, name, shares, account_id)
                   VALUES ($1, $2, $3, (select id from account where email = $4))`
	tx.MustExec(sharesInsert, "sbny", "Signature Bank of NY", 35, "investor@somewhere")
	tx.MustExec(sharesInsert, "aapl", "Apple", 250, "investor@somewhere")
	tx.MustExec(sharesInsert, "brk.a", "Berkshire Hathaway", 12, "investor@somewhere")
	tx.MustExec(sharesInsert, "baba", "Ali Baba Group", 75, "investor@somewhere")
	tx.Commit()
}
