import { Component, OnInit } from '@angular/core';
import { InvoiceFacade } from '../+state/invoice.facade';
import { loadInvoices } from '../+state/invoice.actions';

@Component({
  selector: 'adm-invoice-list-container',
  templateUrl: './invoice-list-container.component.html',
  styleUrls: ['./invoice-list-container.component.scss']
})
export class InvoiceListContainerComponent implements OnInit {

  public readonly invoices$ = this.facade.invoices$;

  constructor(private readonly facade: InvoiceFacade) { }

  ngOnInit() {
    this.facade.dispatch(loadInvoices());
  }

  public openPayWindow(id: string) {
    const height = 700;
    const width = 500;
    const screen = window.screen;
    const left = (screen.width / 2) - (width / 2);
    const top = (screen.height / 2) - (height / 2);
    const feats = `location=no,toolbar=no,width=${width},height=${height},top=${top},left=${left}`;

    // TODO: this URL needs to be a config
    window.open(`${window.location.protocol}//${window.location.hostname}:8080/invoices/${id}`,
      'Pay with VendoPunkto', feats);
  }

}
