package repositories

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/leonardochaia/vendopunkto/errors"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

type postgresInvoiceRepository struct {
	db *pg.DB
}

// NewPostgresInvoiceRepository creates the invoice's postgress implementation
func NewPostgresInvoiceRepository(db *pg.DB) vendopunkto.InvoiceRepository {
	return postgresInvoiceRepository{
		db: db,
	}
}

func (r postgresInvoiceRepository) Search(
	ctx context.Context,
	filter vendopunkto.InvoiceFilter) ([]vendopunkto.Invoice, error) {
	const op errors.Op = "invoiceRepository.search"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	var invoices []vendopunkto.Invoice

	err = tx.Model(&invoices).
		Relation("PaymentMethods").
		Relation("PaymentMethods.Payments").
		Order("created_at DESC").
		Select()
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (r postgresInvoiceRepository) FindByID(ctx context.Context, id string) (*vendopunkto.Invoice, error) {
	const op errors.Op = "invoiceRepository.findByID"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	invoice := &vendopunkto.Invoice{
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

func (r postgresInvoiceRepository) FindByAddress(ctx context.Context, address string) (*vendopunkto.Invoice, error) {
	const op errors.Op = "invoiceRepository.findByAddress"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	method := &vendopunkto.PaymentMethod{}
	err = tx.Model(method).Where("address = ?", address).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.E(op, errors.NotExist, err)
		}
		return nil, errors.E(op, errors.Internal, err)
	}

	return r.FindByID(ctx, method.InvoiceID)
}

func (r postgresInvoiceRepository) Create(ctx context.Context, i *vendopunkto.Invoice) error {
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

func (r postgresInvoiceRepository) CreatePayment(ctx context.Context, payment *vendopunkto.Payment) error {
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

func (r postgresInvoiceRepository) UpdatePayment(ctx context.Context, payment *vendopunkto.Payment) error {
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

func (r postgresInvoiceRepository) UpdatePaymentMethod(ctx context.Context, method *vendopunkto.PaymentMethod) error {
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

func (r postgresInvoiceRepository) GetMaxBlockHeightForCurrencies(ctx context.Context) (map[string]uint64, error) {
	const op errors.Op = "invoiceRepository.findPendingPaymentMethods"
	tx, err := store.GetTransactionFromContextOrCreate(ctx, r.db)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	type PaymentWithCurrency struct {
		vendopunkto.Payment `pg:",inherit"`
		Currency            string
	}

	payments := []PaymentWithCurrency{}
	err = tx.Model(&payments).
		ColumnExpr("COALESCE(MAX(payment.block_height),0) block_height").
		ColumnExpr("LOWER(payment_method.currency) AS currency").
		Join("RIGHT JOIN payment_methods AS payment_method").
		JoinOn("payment.payment_method_id = payment_method.id").
		Group("payment_method.currency").
		Select()

	if err != nil && err != pg.ErrNoRows {
		return nil, errors.E(op, errors.Internal, err)
	}

	output := make(map[string]uint64)
	for _, m := range payments {
		output[m.Currency] = m.BlockHeight
	}

	return output, nil
}
