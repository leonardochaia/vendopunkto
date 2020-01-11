package vendopunkto

import (
	"testing"

	"github.com/leonardochaia/vendopunkto/testutils"
	"github.com/shopspring/decimal"
)

var (
	d2 = decimal.NewFromInt(2)
	d4 = decimal.NewFromInt(2)
)

func TestInvoiceStatus(t *testing.T) {
	total := decimal.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	method := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	if inv.Status() != Pending {
		t.Errorf("Expected invoice to be pending but it's %d", inv.Status())
	}

	payment := method.AddPayment("fake-hash", total, 0, 123)

	if inv.Status() != Pending {
		t.Errorf("Expected invoice to be pending but it's %d", inv.Status())
	}

	payment.Update(10, 123)

	if inv.Status() != Completed {
		t.Errorf("Expected invoice to be completed but it's %d", inv.Status())
	}
}

func TestInvoicePaymentPercentageSameCurrency(t *testing.T) {
	total := decimal.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	method := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	p := inv.CalculatePaymentPercentage()
	if p != 0 {
		t.Errorf("Expected invoice to be 0 but it's %f", p)
	}

	method.AddPayment("fake-hash", total.Div(d2), 1, 123)

	p = inv.CalculatePaymentPercentage()
	if p != 50 {
		t.Errorf("Expected invoice to be 50 but it's %f", p)
	}

	method.AddPayment("fake-hash", total.Div(d4), 1, 123)

	p = inv.CalculatePaymentPercentage()
	if p != 100 {
		t.Errorf("Expected invoice to be 150 but it's %f", p)
	}
}

func TestInvoicePaymentPercentageDiffCurrency(t *testing.T) {
	total := decimal.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	first := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	second := inv.AddPaymentMethod("xmr-double", "xmr2-fake-addr", total.Mul(d2))

	p := inv.CalculatePaymentPercentage()
	if p != 0 {
		t.Errorf("Expected invoice to be 0 but it's %f", p)
	}

	second.AddPayment("fake-hash2", total, 1, 123)

	p = inv.CalculatePaymentPercentage()
	if p != 50 {
		t.Errorf("Expected invoice to be 50 but it's %f", p)
	}

	first.AddPayment("fake-hash", total.Div(d2), 1, 123)

	p = inv.CalculatePaymentPercentage()
	if p != 100 {
		t.Errorf("Expected invoice to be 100 but it's %f", p)
	}
}

func TestInvoiceRemainingAmount(t *testing.T) {
	total := decimal.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	method := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	r := inv.CalculateRemainingAmount()

	testutils.DecimalEquals(t, total, r)

	method.AddPayment("fake-hash", total.Div(d2), 1, 123)

	r = inv.CalculateRemainingAmount()

	testutils.DecimalEquals(t, total.Div(d2), r)

	method.AddPayment("fake-hash", total.Div(d2), 1, 123)

	r = inv.CalculateRemainingAmount()
	testutils.DecimalEquals(t, decimal.Zero, r)
}

func TestInvoicePaymentMethodRemainingAmount(t *testing.T) {
	total := decimal.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	converted := total.Mul(d2)
	method := inv.AddPaymentMethod("xmr-double", "xmr-fake-addr", converted)

	r := inv.CalculatePaymentMethodRemaining(*method)
	testutils.DecimalEquals(t, converted, r)

	method.AddPayment("fake-hash", converted.Div(d2), 1, 123)

	r = inv.CalculatePaymentMethodRemaining(*method)
	testutils.DecimalEquals(t, converted.Div(d2), r)

	method.AddPayment("fake-hash", converted.Div(d2), 1, 123)

	r = inv.CalculatePaymentMethodRemaining(*method)

	testutils.DecimalEquals(t, decimal.Zero, r)
}

func TestInvoicePaymentMethodRemainingAmountOverPayed(t *testing.T) {
	total := decimal.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	method := inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)

	r := inv.CalculateRemainingAmount()

	testutils.DecimalEquals(t, total, r)

	method.AddPayment("fake-hash", total.Mul(d2), 1, 123)

	r = inv.CalculateRemainingAmount()
	testutils.DecimalEquals(t, decimal.Zero, r)
}

func TestDecimalRemaining(t *testing.T) {
	total := decimal.NewFromFloat(0.333)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)
	inv.AddPaymentMethod("xmr2", "xmr-fake-addr", total.Mul(d2)).
		AddPayment("fake-hash", total.Mul(d2), 1, 1)

	expected := decimal.NewFromFloat(0.333)
	result := inv.CalculateTotalPayedAmount()

	testutils.DecimalEquals(t, expected, result)
}

func TestDecimalPayed(t *testing.T) {
	total := decimal.NewFromFloat(0.333)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)
	inv.AddPaymentMethod("xmr2", "xmr-fake-addr", total.Mul(d2)).
		AddPayment("fake-hash", total, 1, 1)

	expected := decimal.NewFromFloat(0.16650)
	result := inv.CalculateTotalPayedAmount()

	testutils.DecimalEquals(t, expected, result)
}

func TestDecimalConversions(t *testing.T) {
	total := decimal.NewFromFloat(1)
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	rate := decimal.NewFromFloat(0.006539)

	converted := total.Mul(rate)

	inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)
	inv.AddPaymentMethod("btc", "btc-fake-addr", converted).
		AddPayment("fake-hash", decimal.NewFromFloat(0.00163475), 1, 1)

	expected := decimal.NewFromFloat(0.25)
	result := inv.CalculateTotalPayedAmount()

	testutils.DecimalEquals(t, expected, result)
}

func TestDecimalPrecission(t *testing.T) {
	total := decimal.NewFromFloat(0.000000000001) // 1 piconero
	inv := &Invoice{
		Currency: "xmr",
		Total:    total,
	}

	rate := decimal.NewFromFloat(0.006539)

	converted := total.Mul(rate)

	t.Logf("c: %s", converted.Div(d2).String())

	inv.AddPaymentMethod("xmr", "xmr-fake-addr", total)
	inv.AddPaymentMethod("btc", "btc-fake-addr", converted).
		AddPayment("fake-hash", converted, 1, 1)

	expected := total
	result := inv.CalculateTotalPayedAmount()

	testutils.DecimalEquals(t, expected, result)
}
