import { Injectable } from '@angular/core';
import { InvoiceState } from './invoice.reducer';
import { Store, Action } from '@ngrx/store';
import {
  selectInvoices,
  selectCurrencies,
  selectLoadingInvoices,
  selectLoadingCurrencies,
  selectLoadingPaymentMethods,
  selectBasicCreationInfo,
  selectCreatingInvoice,
  selectPaymentCurrencies
} from './invoice.selectors';

@Injectable({
  providedIn: 'root'
})
export class InvoiceFacade {

  public readonly invoices$ = this.store.select(selectInvoices);
  public readonly loadingInvoices$ = this.store.select(selectLoadingInvoices);

  public readonly currencies$ = this.store.select(selectCurrencies);
  public readonly paymentCurrencies$ = this.store.select(selectPaymentCurrencies);
  public readonly loadingCurrencies$ = this.store.select(selectLoadingCurrencies);

  public readonly creationBasicInfo$ = this.store.select(selectBasicCreationInfo);
  public readonly loadingPaymentMethods$ = this.store.select(selectLoadingPaymentMethods);
  public readonly creating$ = this.store.select(selectCreatingInvoice);

  constructor(private readonly store: Store<InvoiceState>) { }

  public dispatch(action: Action) {
    this.store.dispatch(action);
  }
}
