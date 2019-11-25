package invoice

import (
	"context"
	"database/sql"

	"github.com/leonardochaia/vendopunkto/store"
	"github.com/leonardochaia/vendopunkto/store/postgres"
)

type Store struct {
	client *postgres.Client
}

func NewStore(client *postgres.Client) (Store, error) {
	return Store{
		client: client,
	}, nil
}

func (s *Store) GetByID(ctx context.Context, id string) (*Invoice, error) {

	b := new(Invoice)
	err := s.client.Database.GetContext(ctx, b, `SELECT * FROM invoices WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return b, nil

}

func (s *Store) SaveInvoice(ctx context.Context, i *Invoice) (string, error) {

	// Generate an ID if needed
	if i.ID == "" {
		i.ID = s.client.NewID()
	}

	_, err := s.client.Database.ExecContext(ctx, `
		INSERT INTO invoices (id, amount, denomination, address)
		VALUES($1, $2, $3, $4)
	`, i.ID, i.Amount, i.Denomination, i.PaymentAddress)

	if err != nil {
		return i.ID, err
	}
	return i.ID, nil

}

// // InvoiceDeleteByID an Invoice
// func (c *Client) InvoiceDeleteByID(ctx context.Context, id string) error {

// 	_, err := c.db.ExecContext(ctx, `DELETE FROM Invoice WHERE id = $1`, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil

// }

// // InvoiceFind gets Invoices
// func (c *Client) InvoiceFind(ctx context.Context) ([]*invoice.Invoice, error) {

// 	var bs = make([]*invoice.Invoice, 0)
// 	err := c.db.SelectContext(ctx, &bs, `SELECT * FROM Invoice`)
// 	if err == sql.ErrNoRows {
// 		// No Error
// 	} else if err != nil {
// 		return bs, err
// 	}
// 	return bs, nil

// }
