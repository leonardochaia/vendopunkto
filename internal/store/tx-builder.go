package store

import "context"

import "github.com/go-pg/pg"

// TransactionBuilder builds database transactions
type TransactionBuilder interface {
	BuildLazyTransactionContext(ctx context.Context) (context.Context, LazyTransaction, error)
}

// NewPostgreTransactionBuilder creates the builder implementation for PostgreSQL
func NewPostgreTransactionBuilder(db *pg.DB) TransactionBuilder {
	return &txBuilderImpl{
		db: db,
	}
}

type txBuilderImpl struct {
	db *pg.DB
}

func (builder *txBuilderImpl) BuildLazyTransactionContext(ctx context.Context) (
	context.Context, LazyTransaction, error) {
	tx := NewLazyTransaction(builder.db)

	ctx = context.WithValue(ctx, TransactionPerRequestKey, tx)
	return ctx, tx, nil
}
