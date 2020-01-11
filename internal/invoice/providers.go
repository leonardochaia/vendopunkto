package invoice

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
)

// Providers for invoice logic
var Providers = wire.NewSet(NewManager, NewTopic)

// NewManager creates an InvoiceManager impl
func NewManager(
	repository vendopunkto.InvoiceRepository,
	pluginManager vendopunkto.PluginManager,
	globalLogger hclog.Logger,
	topic vendopunkto.InvoiceTopic) (vendopunkto.InvoiceManager, error) {
	return &invoiceManager{
		logger:        globalLogger.Named("invoice-manager"),
		pluginManager: pluginManager,
		repository:    repository,
		topic:         topic,
	}, nil
}
