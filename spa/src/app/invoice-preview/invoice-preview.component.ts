import { Component, OnInit, Input } from '@angular/core';
import { switchMap, tap } from 'rxjs/operators';
import { InvoiceDTO, PaymentMethod } from '../model';
import { VendopunktoApiService } from '../vendopunkto-api.service';

@Component({
  selector: 'app-invoice-preview',
  templateUrl: './invoice-preview.component.html',
  styleUrls: ['./invoice-preview.component.scss']
})
export class InvoicePreviewComponent implements OnInit {

  @Input()
  public invoice: InvoiceDTO;

  public paymentMethod: PaymentMethod;

  constructor(private readonly api: VendopunktoApiService) { }

  public ngOnInit() {
    this.paymentMethod = this.getPaymentMethod(this.invoice.currency);
  }

  public verifyPayment() {
    this.api.getInvoice(this.invoice.id)
      .subscribe(freshInvoice => {
        this.invoice = freshInvoice;
      });
  }

  public changePaymentMethod(currency: string) {
    const newMethod = this.getPaymentMethod(currency);
    if (newMethod.address === '') {

      this.api.generatePaymentMethodAddress(this.invoice.id, newMethod.currency)
        .subscribe(invoice => {
          this.invoice = invoice;
          this.paymentMethod = this.getPaymentMethod(currency);
        });
    } else {
      this.paymentMethod = newMethod;
    }

  }

  protected getPaymentMethod(currency: string) {
    return this.invoice.paymentMethods.filter(pm => pm.currency === currency)[0];
  }

}
