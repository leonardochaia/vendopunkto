package invoice

import "context"

type InvoiceFilter struct {
}

// InvoiceRepository is the abstraction for handling Invoices database
type InvoiceRepository interface {
	Search(ctx context.Context, filter InvoiceFilter) ([]Invoice, error)
	FindByID(ctx context.Context, id string) (*Invoice, error)
	FindByAddress(ctx context.Context, address string) (*Invoice, error)
	Create(ctx context.Context, invoice *Invoice) error
	UpdatePaymentMethod(ctx context.Context, method *PaymentMethod) error
	CreatePayment(ctx context.Context, payment *Payment) error
	UpdatePayment(ctx context.Context, payment *Payment) error

	// GetMaxBlockHeightForCurrencies returns a map of currencies and the last
	// known block height. It must return all currencies that are awaiting
	// payment, even if they don't have a known last block height.
	GetMaxBlockHeightForCurrencies(ctx context.Context) (map[string]uint64, error)
}
