package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fatih/color"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"os"
)

type DbConfig struct {
	DatabaseName   string
	TursoAuthToken string
}

func Init(databaseName, tursoAuthToken string) (DbConfig, error) {

	color.Yellow("initialising database connection configuration...")

	if len(databaseName) == 0 || len(tursoAuthToken) == 0 {
		return DbConfig{}, errors.New("no database name or auth token provided")
	}

	color.Green("successfully configured the connection!")
	return DbConfig{
		DatabaseName:   databaseName,
		TursoAuthToken: tursoAuthToken,
	}, nil
}

func (database DbConfig) CreateConnection() *sql.DB {

	color.Yellow("starting connection...")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s",
		database.DatabaseName,
		database.TursoAuthToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	color.Green("successfully connected to database: " + database.DatabaseName)

	// check when to close database connection
	// defer db.Close()
	return db
}

func (database DbConfig) CloseConnection(db *sql.DB) {
	defer db.Close()
}

// test func for now to query users

type User struct {
	ID        string
	Username  string
	Password  string
	createdOn string
	lastLogin string
}

// broken func just here for reference for now
func QueryUsers(database *sql.DB) []User {

	rows, err := database.Query("SELECT * FROM users")
	if err != nil {
		color.Red("failed to execute query %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			color.Red("error scanning row: %v", err)
			return []User{}
		}

		users = append(users, user)
		color.Green(string(user.ID), user.Name)
	}

	if err := rows.Err(); err != nil {
		color.Red("error during row iteration: %v", err)
	}

	return users
}
