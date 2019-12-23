package store

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/go-pg/pg"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/conf"
)

// NewDB Initialices the database
func NewDB(
	globalLogger hclog.Logger,
	startupConf conf.Startup,
) (*pg.DB, error) {
	var (
		username string = startupConf.Storage.Username
		password string = startupConf.Storage.Password
		hostname string = startupConf.Storage.Host
		port     string = startupConf.Storage.Port
		database string = startupConf.Storage.Database
		sslMode  string = startupConf.Storage.SSLMode
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

	db := pg.Connect(&pg.Options{
		User:     username,
		Password: password,
		Database: "postgres",
		Addr:     net.JoinHostPort(hostname, port),
	})

	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	db.AddQueryHook(dbLogger{
		logger: globalLogger.Named("db"),
	})

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

type dbLogger struct {
	logger hclog.Logger
}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	query, _ := q.FormattedQuery()
	d.logger.Debug("DB", "query", query)
}
