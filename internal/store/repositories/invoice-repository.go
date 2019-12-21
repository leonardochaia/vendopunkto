package repositories

import (
	"context"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/internal/invoice"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

type postgresInvoiceRepository struct {
	db *pg.DB
}

// NewPostgresInvoiceRepository creates the invoice's postgress implementation
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
	const op errors.Op = "invoiceRepository.findByID"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
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
		if err == pg.ErrNoRows {
			return nil, errors.E(op, errors.NotExist, err)
		}
		return nil, errors.E(op, errors.Internal, err)
	}

	return invoice, nil
}

func (r postgresInvoiceRepository) FindByAddress(ctx context.Context, address string) (*invoice.Invoice, error) {
	const op errors.Op = "invoiceRepository.findByAddress"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	method := &invoice.PaymentMethod{}
	err = tx.Model(method).Where("address = ?", address).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.E(op, errors.NotExist, err)
		}
		return nil, errors.E(op, errors.Internal, err)
	}

	return r.FindByID(ctx, method.InvoiceID)
}

func (r postgresInvoiceRepository) Create(ctx context.Context, i *invoice.Invoice) error {
	const op errors.Op = "invoiceRepository.create"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	_, err = tx.Model(i).Returning("*").Insert(i)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	for _, method := range i.PaymentMethods {
		method.InvoiceID = i.ID
		_, err := tx.Model(method).Returning("*").Insert(method)
		if err != nil {
			return errors.E(op, errors.Internal, err)
		}

		for _, payment := range method.Payments {
			payment.PaymentMethodID = method.ID
			_, err := tx.Model(payment).Returning("*").Insert(payment)
			if err != nil {
				return errors.E(op, errors.Internal, err)
			}
		}
	}

	return nil
}

func (r postgresInvoiceRepository) CreatePayment(ctx context.Context, payment *invoice.Payment) error {
	const op errors.Op = "invoiceRepository.createPayment"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	err = tx.Insert(payment)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	return nil
}

func (r postgresInvoiceRepository) UpdatePayment(ctx context.Context, payment *invoice.Payment) error {
	const op errors.Op = "invoiceRepository.updatePayment"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	err = tx.Update(payment)

	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	return nil
}

func (r postgresInvoiceRepository) UpdatePaymentMethod(ctx context.Context, method *invoice.PaymentMethod) error {
	const op errors.Op = "invoiceRepository.updatePaymentMethod"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	err = tx.Update(method)

	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	return nil
}

func createSchema(db *pg.DB) error {
	const op errors.Op = "invoiceRepository.createSchema"
	for _, model := range []interface{}{(*invoice.Invoice)(nil), (*invoice.PaymentMethod)(nil), (*invoice.Payment)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return errors.E(op, errors.Internal, err)
		}
	}
	return nil
}
