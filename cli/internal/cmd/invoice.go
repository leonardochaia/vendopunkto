package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/unit"
	"github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
)

func init() {

	invoiceViewCmd.Flags().BoolP("qr", "q", false, "Show QR code")
	invoiceViewCmd.Flags().StringP("method", "m", "", "Payment method. i.e --method=xmr")
	invoiceCmd.AddCommand(invoiceCreateCmd)
	invoiceCmd.AddCommand(invoiceViewCmd)

	rootCmd.AddCommand(invoiceCmd)
}

var (
	qrTerminalConfig = qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
	}

	invoiceCmd = &cobra.Command{
		Use:   "invoices",
		Short: "Manage Invoices",
		Long:  `Manage Invoices`,
	}

	invoiceCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an Invoices",
		Long:  "Create an Invoices",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error { // Initialize the databse

			currency := args[0]
			total, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return err
			}

			invoice, err := publicClient.CreateInvoice(unit.NewFromFloat(total), currency)

			if err != nil {
				return err
			}

			fmt.Print("Your invoice has been created and is awaiting payment.\n")
			printInvoice(invoice)

			return nil
		},
	}

	invoiceViewCmd = &cobra.Command{
		Use:   "view",
		Short: "View an Invoice",
		Long:  "View an Invoice",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error { // Initialize the databse

			ID := args[0]

			qr, err := cmd.Flags().GetBool("qr")
			if err != nil {
				return err
			}

			method, err := cmd.Flags().GetString("method")
			if err != nil {
				return err
			}

			invoice, err := publicClient.GetInvoice(ID)

			if err != nil {
				return err
			}

			if method == "" {
				method = invoice.Currency
			} else {
				invoice, err = ensureMethodAddressExists(invoice, method)
				if err != nil {
					return err
				}
			}

			fmt.Print("Found invoice\n")
			printInvoice(invoice)

			if qr {
				for _, m := range invoice.PaymentMethods {
					if m.Currency == method {
						printQR(m)
						break
					}
				}
			}

			return nil
		},
	}
)

func printInvoice(invoice *dtos.InvoiceDto) {
	fmt.Printf("  Invoice ID: %s\n", invoice.ID)
	fmt.Printf("  Total: %s %s\n", strings.ToUpper(invoice.Currency), invoice.Total.ValueFormatted)
	fmt.Printf("  Status: %s %s%%\n", getInvoiceStatus(invoice.Status),
		strconv.FormatFloat(invoice.PaymentPercentage, 'f', -1, 64))
	if invoice.Total.Value != invoice.Remaining.Value {
		fmt.Printf("  Remaining: %s %s\n", strings.ToUpper(invoice.Currency), invoice.Remaining.ValueFormatted)
	}

	fmt.Print("You can pay using any of the following methods:\n")

	for _, method := range invoice.PaymentMethods {
		fmt.Printf("  %s %s %s\n",
			strings.ToUpper(method.Currency),
			method.Remaining.ValueFormatted,
			method.Address)
	}
}

func ensureMethodAddressExists(invoice *dtos.InvoiceDto, currency string) (*dtos.InvoiceDto, error) {
	for _, method := range invoice.PaymentMethods {
		if method.Currency == currency {

			if method.Address == "" {
				freshInv, err := publicClient.GeneratePaymentMethodAdress(invoice.ID, method.Currency)
				if err != nil {
					return nil, err
				}
				return freshInv, nil
			}
			break
		}

	}

	return invoice, nil
}

func printQR(method *dtos.PaymentMethodDto) {
	if method.QRCode == "" {
		return
	}

	fmt.Printf("\nPrinting QR Code for %s\n", strings.ToUpper(method.Currency))
	qrterminal.GenerateWithConfig(method.QRCode, qrTerminalConfig)
	fmt.Printf("%s %s\n", strings.ToUpper(method.Currency), method.Remaining.ValueFormatted)
	fmt.Printf("Address: %s\n", method.Address)
}

func getInvoiceStatus(status uint) string {
	switch status {
	case 1:
		return "Pending"
	case 2:
		return "Completed"
	case 3:
		return "Failed"
	}
	return "unknown"
}
