package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fatih/color"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"os"
)

type Db struct {
	DatabaseName   string
	TursoAuthToken string
}

func Init(databaseName, tursoAuthToken string) (Db, error) {

	color.Yellow("initialising database connection configuration...")

	if len(databaseName) == 0 || len(tursoAuthToken) == 0 {
		return Db{}, errors.New("no database name or auth token provided")
	}

	color.Green("successfully configured the connection!")
	return Db{
		DatabaseName:   databaseName,
		TursoAuthToken: tursoAuthToken,
	}, nil
}

func (database Db) CreateConnection() {

	color.Yellow("starting connection...")

	url := fmt.Sprintf("libsql://%s.turso.io=authToken=%s",
		database.DatabaseName,
		database.TursoAuthToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	color.Green("successfully connected to database: " + database.DatabaseName)

	defer db.Close()
}
