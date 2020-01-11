package store

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/leonardochaia/vendopunkto/errors"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
)

// TODO: Migrations are not really a thing.
// we just start from scratch and create every table.

var entities = []interface{}{
	(*vendopunkto.Currency)(nil),
	(*vendopunkto.Invoice)(nil),
	(*vendopunkto.PaymentMethod)(nil),
	(*vendopunkto.Payment)(nil),
}

// doMigrations performs Postgres database migrations
func doMigrations(db *pg.DB) error {
	const op errors.Op = "repositories.doMigrations"
	for _, model := range entities {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return errors.E(op, errors.Internal, err)
		}
	}
	return nil
}
