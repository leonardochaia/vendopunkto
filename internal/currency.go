package vendopunkto

import "context"

// Currency holds metada for a currency.
type Currency struct {
	Symbol       string `sql:",pk"`
	Name         string `sql:",notnull"`
	LogoImageURL string
}

// CurrencyRepository is the abstraction for handling currencies database
type CurrencyRepository interface {
	Search(ctx context.Context) ([]Currency, error)
	FindBySymbol(ctx context.Context, symbol string) (*Currency, error)
	SelectOrInsert(ctx context.Context, currency *Currency) (*Currency, error)
	AddOrUpdate(ctx context.Context, currency *Currency) (*Currency, error)
}
