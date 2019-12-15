package invoice

import (
	"github.com/leonardochaia/vendopunkto/unit"
	"testing"
)

func TestInvoiceStatus(t *testing.T) {
	total := unit.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	method := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	if inv.Status() != Pending {
		t.Errorf("Expected invoice to be pending but it's %d", inv.Status())
	}

	payment := method.AddPayment("fake-hash", total, 0)

	if inv.Status() != Pending {
		t.Errorf("Expected invoice to be pending but it's %d", inv.Status())
	}

	payment.Update(10)

	if inv.Status() != Completed {
		t.Errorf("Expected invoice to be completed but it's %d", inv.Status())
	}
}

func TestInvoicePaymentPercentageSameCurrency(t *testing.T) {
	total := unit.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	method := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	p := inv.CalculatePaymentPercentage()
	if p != 0 {
		t.Errorf("Expected invoice to be 0 but it's %f", p)
	}

	method.AddPayment("fake-hash", total/2, 1)

	p = inv.CalculatePaymentPercentage()
	if p != 50 {
		t.Errorf("Expected invoice to be 50 but it's %f", p)
	}

	method.AddPayment("fake-hash", total/4, 1)
	method.AddPayment("fake-hash", total/4, 1)

	p = inv.CalculatePaymentPercentage()
	if p != 100 {
		t.Errorf("Expected invoice to be 50 but it's %f", p)
	}
}

func TestInvoicePaymentPercentageDiffCurrency(t *testing.T) {
	total := unit.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	first := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	second := inv.AddPaymentMethod("xmr-double", "xmr2-fake-addr", total*2)

	p := inv.CalculatePaymentPercentage()
	if p != 0 {
		t.Errorf("Expected invoice to be 0 but it's %f", p)
	}

	second.AddPayment("fake-hash2", total, 1)

	p = inv.CalculatePaymentPercentage()
	if p != 50 {
		t.Errorf("Expected invoice to be 50 but it's %f", p)
	}

	first.AddPayment("fake-hash", total/2, 1)

	p = inv.CalculatePaymentPercentage()
	if p != 100 {
		t.Errorf("Expected invoice to be 100 but it's %f", p)
	}
}

func TestInvoiceRemainingAmount(t *testing.T) {
	total := unit.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	method := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	r := inv.CalculateRemainingAmount()
	if r != total {
		t.Errorf("Expected invoice to be %d but it's %d", total, r)
	}

	method.AddPayment("fake-hash", total/2, 1)

	r = inv.CalculateRemainingAmount()
	if r != total/2 {
		t.Errorf("Expected invoice to be %d but it's %d", total/2, r)
	}

	method.AddPayment("fake-hash", total/2, 1)

	r = inv.CalculateRemainingAmount()
	if r != 0 {
		t.Errorf("Expected invoice to be 0 but it's %d", r)
	}
}

func TestInvoicePaymentMethodRemainingAmount(t *testing.T) {
	total := unit.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	converted := total * 2
	method := inv.AddPaymentMethod("xmr-double", "xmr-fake-addr", converted)

	r := inv.CalculatePaymentMethodRemaining(*method)
	if r != converted {
		t.Errorf("Expected invoice to be %d but it's %d", converted, r)
	}

	method.AddPayment("fake-hash", converted/2, 1)

	r = inv.CalculatePaymentMethodRemaining(*method)
	if r != converted/2 {
		t.Errorf("Expected invoice to be %d but it's %d", converted/2, r)
	}

	method.AddPayment("fake-hash", converted/2, 1)

	r = inv.CalculatePaymentMethodRemaining(*method)
	if r != 0 {
		t.Errorf("Expected invoice to be 0 but it's %d", r)
	}
}

func TestInvoicePaymentMethodRemainingAmountOverPayed(t *testing.T) {
	total := unit.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	method := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	r := inv.CalculateRemainingAmount()
	if r != total {
		t.Errorf("Expected invoice to be %d but it's %d", total, r)
	}

	method.AddPayment("fake-hash", total*2, 1)

	r = inv.CalculateRemainingAmount()
	if r != 0 {
		t.Errorf("Expected invoice to be 0 but it's %d", r)
	}
}
