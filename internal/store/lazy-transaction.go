package store

import (
	"sync"

	"github.com/go-pg/pg"
)

type LazyTransaction interface {
	HasTransaction() bool
	GetCurrent() (*pg.Tx, error)
	CommitIfNeeded() error
	RollbackIfNeeded() error
}

func NewLazyTransaction(db *pg.DB) LazyTransaction {
	return &lazyTransaction{
		db:   db,
		once: sync.Once{},
	}
}

type lazyTransaction struct {
	db        *pg.DB
	once      sync.Once
	currentTx *pg.Tx
}

func (lt *lazyTransaction) GetCurrent() (*pg.Tx, error) {
	var (
		err error
	)

	lt.once.Do(func() {
		lt.currentTx, err = lt.db.Begin()
	})

	return lt.currentTx, err
}

func (lt *lazyTransaction) CommitIfNeeded() error {
	if lt.currentTx != nil {
		return lt.currentTx.Commit()
	}
	return nil
}

func (lt *lazyTransaction) RollbackIfNeeded() error {
	if lt.currentTx != nil {
		return lt.currentTx.Rollback()
	}
	return nil
}

func (lt *lazyTransaction) HasTransaction() bool {
	return lt.currentTx != nil
}
