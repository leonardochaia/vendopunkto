import { Injectable } from '@angular/core';
import { InvoiceState } from './invoice.reducer';
import { Store, Action } from '@ngrx/store';
import { selectInvoices, selectCurrencies } from './invoice.selectors';

@Injectable({
  providedIn: 'root'
})
export class InvoiceFacade {

  public readonly invoices$ = this.store.select(selectInvoices);
  public readonly currencies$ = this.store.select(selectCurrencies);

  constructor(private readonly store: Store<InvoiceState>) { }

  public dispatch(action: Action) {
    this.store.dispatch(action);
  }
}
