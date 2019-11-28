package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func NewDB() (*gorm.DB, error) {
	var (
		username string = viper.GetString("storage.username")
		password string = viper.GetString("storage.password")
		hostname string = viper.GetString("storage.host")
		port     string = viper.GetString("storage.port")
		database string = viper.GetString("storage.database")
		sslMode  string = viper.GetString("storage.ssl_mode")
	)

	// Username
	if username == "" {
		return nil, fmt.Errorf("No storage username specified")
	}

	// Password
	if password == "" {
		return nil, fmt.Errorf("No storage password specified")
	}

	// Host
	if hostname == "" {
		return nil, fmt.Errorf("No storage hostname specified")
	}

	// Port
	if port == "" {
		return nil, fmt.Errorf("No storage port specified")
	}

	// DB Name
	if database == "" {
		return nil, fmt.Errorf("No storage database name specified")
	}

	// SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dbURLOptions := fmt.Sprintf("?sslmode=%s", sslMode)

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s",
		username, password, hostname, port)

	err := createDatabase(dbURI+dbURLOptions, database)

	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("postgres", dbURI+"/"+database+dbURLOptions)
	// db.DB().SetMaxOpenConns(maxConnections)

	return db, err
}

func createDatabase(url string, dbName string) error {
	createDb, err := sql.Open("postgres", url)
	// Attempt to create the database if it doesn't exist
	if err != nil {
		return err
	}

	defer createDb.Close()

	var one sql.NullInt64
	err = createDb.QueryRow(`SELECT 1 from pg_database WHERE datname=$1`, dbName).Scan(&one)
	if err == nil {
		// already exists
		return nil
	} else if err != sql.ErrNoRows && !strings.Contains(err.Error(), "does not exist") {
		// Some other error besides does not exist
		return fmt.Errorf("Could not check for database: %s", err)
	}

	_, err = createDb.Exec(`CREATE DATABASE ` + dbName)
	if err != nil {
		return fmt.Errorf("Could not create database: %s", err)
	}

	return nil
}
