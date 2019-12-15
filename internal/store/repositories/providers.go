package repositories

import "github.com/google/wire"

var Providers = wire.NewSet(NewPostgresInvoiceRepository)
