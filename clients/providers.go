package clients

import (
	"github.com/google/wire"
)

// Providers to use with wire
var Providers = wire.NewSet(NewHTTPClient)
