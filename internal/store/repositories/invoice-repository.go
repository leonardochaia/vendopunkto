package repositories

import "github.com/jinzhu/gorm"

import "github.com/leonardochaia/vendopunkto/internal/invoice"

type postgresInvoiceRepository struct {
	db *gorm.DB
}

func NewPostgresInvoiceRepository(db *gorm.DB) invoice.InvoiceRepository {

	// TODO: Migrations
	db.AutoMigrate(&invoice.Invoice{})
	db.AutoMigrate(&invoice.Payment{})
	db.AutoMigrate(&invoice.PaymentMethod{})

	return postgresInvoiceRepository{
		db: db,
	}
}

func (r postgresInvoiceRepository) FindByID(id string) (*invoice.Invoice, error) {
	return r.findInvoice("ID", id)
}

func (r postgresInvoiceRepository) FindByAddress(address string) (*invoice.Invoice, error) {
	var po invoice.PaymentMethod
	result := r.db.Where("address = ?", address).First((&po))

	if result.RecordNotFound() {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return r.findInvoice("ID", po.InvoiceID)
}

func (r postgresInvoiceRepository) Create(invoice *invoice.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r postgresInvoiceRepository) CreatePayment(payment *invoice.Payment) error {
	return r.db.Create(payment).Error
}

func (r postgresInvoiceRepository) UpdatePayment(payment *invoice.Payment) error {
	return r.db.Save(payment).Error
}

func (r postgresInvoiceRepository) UpdatePaymentMethod(method *invoice.PaymentMethod) error {
	return r.db.Save(method).Error
}

func (r *postgresInvoiceRepository) findInvoice(key string, value string) (*invoice.Invoice, error) {
	var invoice invoice.Invoice

	result := r.db.
		Preload("PaymentMethods").
		Preload("PaymentMethods.Payments").
		First(&invoice, key+" = ?", value)

	if result.RecordNotFound() {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &invoice, nil
}
