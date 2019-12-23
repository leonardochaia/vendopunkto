import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { InvoiceDTO } from '../model';
import { VendopunktoApiService } from '../vendopunkto-api.service';

@Component({
  selector: 'app-invoice-search-container',
  templateUrl: './invoice-search-container.component.html',
  styleUrls: ['./invoice-search-container.component.scss']
})
export class InvoiceSearchContainerComponent implements OnInit {

  constructor(
    private readonly snapshot: ActivatedRoute,
    private readonly api: VendopunktoApiService,
  ) { }

  public current: InvoiceDTO;

  ngOnInit() {
    this.snapshot.paramMap
      .subscribe(params => {
        const invoiceID = params.get('invoiceID');
        if (invoiceID) {
          this.api.getInvoice(invoiceID)
            .subscribe(invoice => {
              this.current = invoice;
            });
        }
      });
  }

}
