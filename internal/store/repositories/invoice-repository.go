package repositories

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
)

type postgresInvoiceRepository struct {
	db *pg.DB
}

func NewPostgresInvoiceRepository(db *pg.DB) (invoice.InvoiceRepository, error) {

	// TODO: Migrations
	err := createSchema(db)
	if err != nil {
		return nil, err
	}

	return postgresInvoiceRepository{
		db: db,
	}, err
}

func (r postgresInvoiceRepository) FindByID(id string) (*invoice.Invoice, error) {
	invoice := &invoice.Invoice{
		ID: id,
	}

	err := r.db.Model(invoice).
		Column("invoice.*").
		Relation("PaymentMethods").
		Relation("PaymentMethods.Payments").
		WherePK().
		First()

	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (r postgresInvoiceRepository) FindByAddress(address string) (*invoice.Invoice, error) {

	method := &invoice.PaymentMethod{}
	err := r.db.Model(method).Where("address = ?", address).Select()
	if err != nil {
		return nil, err
	}

	return r.FindByID(method.InvoiceID)
}

func (r postgresInvoiceRepository) Create(i *invoice.Invoice) error {
	_, err := r.db.Model(i).Returning("*").Insert(i)

	for _, method := range i.PaymentMethods {
		method.InvoiceID = i.ID
		_, err := r.db.Model(method).Returning("*").Insert(method)
		if err != nil {
			return err
		}

		for _, payment := range method.Payments {
			payment.PaymentMethodID = method.ID
			_, err := r.db.Model(payment).Returning("*").Insert(payment)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (r postgresInvoiceRepository) CreatePayment(payment *invoice.Payment) error {
	return r.db.Insert(payment)
}

func (r postgresInvoiceRepository) UpdatePayment(payment *invoice.Payment) error {
	return r.db.Update(payment)
}

func (r postgresInvoiceRepository) UpdatePaymentMethod(method *invoice.PaymentMethod) error {
	return r.db.Update(method)
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*invoice.Invoice)(nil), (*invoice.PaymentMethod)(nil), (*invoice.Payment)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
