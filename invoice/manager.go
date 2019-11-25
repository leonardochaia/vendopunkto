package invoice

import (
	"context"

	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
)

type Manager struct {
	wallet wallet.Client
	store  Store
}

func NewManager(store Store, wallet wallet.Client) (*Manager, error) {
	return &Manager{
		wallet: wallet,
		store:  store,
	}, nil
}

func (inv *Manager) GetInvoice(ctx context.Context, id string) (*Invoice, error) {
	return inv.store.GetByID(ctx, id)
}

func (inv *Manager) CreateInvoice(ctx context.Context, amount uint, denomination string) (*Invoice, error) {

	address, err := inv.wallet.CreateAddress(&wallet.RequestCreateAddress{})
	if err != nil {
		return nil, err
	}

	invoice := &Invoice{
		PaymentAddress: address.Address,
		Amount:         amount,
		Denomination:   denomination,
	}

	_, err = inv.store.SaveInvoice(ctx, invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}
