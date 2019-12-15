package repositories

import (
	"context"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/store"
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

func (r postgresInvoiceRepository) FindByID(ctx context.Context, id string) (*invoice.Invoice, error) {
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, err
	}

	invoice := &invoice.Invoice{
		ID: id,
	}

	err = tx.Model(invoice).
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

func (r postgresInvoiceRepository) FindByAddress(ctx context.Context, address string) (*invoice.Invoice, error) {
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, err
	}

	method := &invoice.PaymentMethod{}
	err = tx.Model(method).Where("address = ?", address).Select()
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, method.InvoiceID)
}

func (r postgresInvoiceRepository) Create(ctx context.Context, i *invoice.Invoice) error {
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return err
	}

	_, err = tx.Model(i).Returning("*").Insert(i)
	if err != nil {
		return err
	}

	for _, method := range i.PaymentMethods {
		method.InvoiceID = i.ID
		_, err := tx.Model(method).Returning("*").Insert(method)
		if err != nil {
			return err
		}

		for _, payment := range method.Payments {
			payment.PaymentMethodID = method.ID
			_, err := tx.Model(payment).Returning("*").Insert(payment)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r postgresInvoiceRepository) CreatePayment(ctx context.Context, payment *invoice.Payment) error {
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return err
	}
	return tx.Insert(payment)
}

func (r postgresInvoiceRepository) UpdatePayment(ctx context.Context, payment *invoice.Payment) error {
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return err
	}
	return tx.Update(payment)
}

func (r postgresInvoiceRepository) UpdatePaymentMethod(ctx context.Context, method *invoice.PaymentMethod) error {
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return err
	}
	return tx.Update(method)
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
