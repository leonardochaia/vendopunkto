package invoice

import "context"

type InvoiceRepository interface {
	FindByID(ctx context.Context, id string) (*Invoice, error)
	FindByAddress(ctx context.Context, address string) (*Invoice, error)
	Create(ctx context.Context, invoice *Invoice) error
	UpdatePaymentMethod(ctx context.Context, method *PaymentMethod) error
	CreatePayment(ctx context.Context, payment *Payment) error
	UpdatePayment(ctx context.Context, payment *Payment) error
}
