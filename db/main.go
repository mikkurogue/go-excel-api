package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Db struct {
	DatabaseName   string
	TursoAuthToken string
}

func Init(databaseName, tursoAuthToken string) (Db, error) {
	if len(databaseName) == 0 || len(tursoAuthToken) == 0 {
		return Db{}, errors.New("no database name or auth token provided")
	}

	return Db{
		DatabaseName:   databaseName,
		TursoAuthToken: tursoAuthToken,
	}, nil
}

func (database Db) CreateConnection() {

	url := fmt.Sprintf("libsql://%s.turso.io=authToken=%s", database.DatabaseName, database.TursoAuthToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	defer db.Close()
}
