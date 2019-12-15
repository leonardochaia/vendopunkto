package invoice

type InvoiceRepository interface {
	FindByID(id string) (*Invoice, error)
	FindByAddress(address string) (*Invoice, error)
	Create(invoice *Invoice) error
	UpdatePaymentMethod(method *PaymentMethod) error
	CreatePayment(payment *Payment) error
	UpdatePayment(payment *Payment) error
}
