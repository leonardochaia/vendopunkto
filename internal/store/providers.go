package store

import "github.com/google/wire"

// Providers are the wire providers for store
var Providers = wire.NewSet(NewDB, NewPostgreTransactionBuilder)
