import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';

@Component({
  selector: 'app-invoice-search-container',
  templateUrl: './invoice-search-container.component.html',
  styleUrls: ['./invoice-search-container.component.scss']
})
export class InvoiceSearchContainerComponent {

  constructor(
    private readonly snapshot: ActivatedRoute,
  ) { }

  public get invoiceId$() {
    return this.snapshot.paramMap
      .pipe(
        map(params => params.get('invoiceID'))
      );
  }
}
