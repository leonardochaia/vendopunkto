package repositories

import "github.com/google/wire"

// Providers for Postgres
var Providers = wire.NewSet(
	NewPostgresInvoiceRepository,
	NewPostgresCurrencyRepository,
)
