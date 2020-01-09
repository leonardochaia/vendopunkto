import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType, OnInitEffects } from '@ngrx/effects';
import { catchError, map, concatMap, tap, withLatestFrom } from 'rxjs/operators';
import { EMPTY, of } from 'rxjs';

import * as InvoiceActions from './invoice.actions';
import { VendopunktoApiService } from '../../vendopunkto-api.service';
import { Action } from '@ngrx/store';
import { Router } from '@angular/router';
import { InvoiceFacade } from './invoice.facade';
import { PaymentMethodCreationParams } from 'shared/shared';

@Injectable()
export class InvoiceEffects implements OnInitEffects {


  loadInvoices$ = createEffect(() => {
    return this.actions$.pipe(

      ofType(InvoiceActions.loadInvoices),
      concatMap(() =>
        this.api.searchInvoices({}).pipe(
          map(data => InvoiceActions.loadInvoicesSuccess({ invoices: data })),
          catchError(error => of(InvoiceActions.loadInvoicesFailure({ error: error.message }))))
      )
    );
  });

  loadCurrencies$ = createEffect(() => {
    return this.actions$.pipe(

      ofType(InvoiceActions.loadCurrencies),
      concatMap(() =>
        this.api.getCurrencies().pipe(
          map(data => InvoiceActions.loadCurrenciesSuccess({
            currencies: data.reduce((a, b) => (a[b.symbol] = b, a), {})
          })),
          catchError(error => of(InvoiceActions.loadCurrenciesFailure({ error: error.message }))))
      )
    );
  });

  createInvoice$ = createEffect(() => {
    return this.actions$.pipe(

      ofType(InvoiceActions.startCreateInvoice),
      concatMap((action) =>
        this.api.createInvoice(action.invoice).pipe(
          map(data => InvoiceActions.createInvoiceSuccess({ invoice: data })),
          tap(() => this.router.navigate(['/invoices'])),
          catchError(error => of(InvoiceActions.createInvoiceFailure({ error: error.message }))))
      )
    );
  });

  invoiceFormChanged$ = createEffect(() => {
    return this.actions$.pipe(

      ofType(InvoiceActions.invoiceCreationFormChanged),
      withLatestFrom(this.facade.creationBasicInfo$, this.facade.paymentCurrencies$),
      concatMap(([action, prevInfo, currencies]) => {
        let info = action.form;
        if (!info.paymentMethods || info.paymentMethods.length === 0) {
          info = {
            ...info,
            paymentMethods: currencies.map(c => ({
              currency: c.symbol,
              total: null,
            } as PaymentMethodCreationParams))
          };
        }
        return this.api.getCurrencyExchange({
          amount: info.total,
          fromCurrency: info.currency,
          toCurrencies: info.paymentMethods.map(pm => pm.currency),
        }).pipe(
          map(result => InvoiceActions.getPaymentMethodExchangeRateSuccess({ result })),
          catchError(error => of(InvoiceActions.getPaymentMethodExchangeRateFailure({ error: error.message })))
        )
      })
    );
  });


  ngrxOnInitEffects(): Action {
    return InvoiceActions.loadCurrencies();
  }

  constructor(
    private readonly actions$: Actions,
    private readonly facade: InvoiceFacade,
    private readonly api: VendopunktoApiService,
    private readonly router: Router) { }

}
