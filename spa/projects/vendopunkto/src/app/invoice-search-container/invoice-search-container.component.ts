import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { VendopunktoApiService } from '../vendopunkto-api.service';
import { map } from 'rxjs/operators';

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

  public get invoiceId$() {
    return this.snapshot.paramMap
      .pipe(
        map(params => params.get('invoiceID'))
      );
  }

  ngOnInit() {
  }

}
