package repositories

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/leonardochaia/vendopunkto/errors"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

type postgresCurrencyRepository struct {
	db *pg.DB
}

// NewPostgresCurrencyRepository creates the invoice's postgress implementation
func NewPostgresCurrencyRepository(db *pg.DB) vendopunkto.CurrencyRepository {
	return postgresCurrencyRepository{
		db: db,
	}
}

func (r postgresCurrencyRepository) Search(
	ctx context.Context) ([]vendopunkto.Currency, error) {
	const op errors.Op = "currencyRepository.search"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	var currencies []vendopunkto.Currency

	err = tx.Model(&currencies).
		Order("symbol ASC").
		Select()
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

func (r postgresCurrencyRepository) FindBySymbol(ctx context.Context, symbol string) (*vendopunkto.Currency, error) {
	const op errors.Op = "currencyRepository.findBySymbol"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	c := &vendopunkto.Currency{
		Symbol: symbol,
	}

	err = tx.Model(c).
		WherePK().
		First()

	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.E(op, errors.NotExist, err)
		}
		return nil, errors.E(op, errors.Internal, err)
	}

	return c, nil
}

func (r postgresCurrencyRepository) SelectOrInsert(
	ctx context.Context,
	c *vendopunkto.Currency) (*vendopunkto.Currency, error) {

	const op errors.Op = "currencyRepository.selectOrInsert"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)

	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	_, err = tx.Model(c).
		WherePK().
		OnConflict("DO NOTHING").
		SelectOrInsert()

	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	return c, nil
}
