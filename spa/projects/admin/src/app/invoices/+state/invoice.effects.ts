import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, concatMap, tap, withLatestFrom } from 'rxjs/operators';
import { of } from 'rxjs';

import * as InvoiceActions from './invoice.actions';
import { VendopunktoApiService } from '../../vendopunkto-api.service';
import { Router } from '@angular/router';
import { PaymentMethodCreationParams } from 'shared/shared';
import { CurrenciesFacade } from '../../currencies/+state/currencies.facade';

@Injectable()
export class InvoiceEffects {


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
      withLatestFrom(this.currenciesFacade.paymentCurrencies$),
      concatMap(([action, paymentCurrencies]) => {
        const info = {
          ...action.form,
          paymentMethods: Object.keys(paymentCurrencies)
            .map(k => {
              const c = paymentCurrencies[k];
              return {
                currency: c.symbol,
                total: null,
              } as PaymentMethodCreationParams;
            })
        };
        return this.api.getCurrencyExchange({
          amount: info.total,
          fromCurrency: info.currency,
          toCurrencies: info.paymentMethods.map(pm => pm.currency),
        }).pipe(
          map(result => InvoiceActions.getPaymentMethodExchangeRateSuccess({ result })),
          catchError(error => of(InvoiceActions.getPaymentMethodExchangeRateFailure({ error: error.message })))
        );
      })
    );
  });

  constructor(
    private readonly actions$: Actions,
    private readonly currenciesFacade: CurrenciesFacade,
    private readonly api: VendopunktoApiService,
    private readonly router: Router) { }

}
