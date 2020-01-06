import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType, OnInitEffects } from '@ngrx/effects';
import { catchError, map, concatMap, tap } from 'rxjs/operators';
import { EMPTY, of } from 'rxjs';

import * as InvoiceActions from './invoice.actions';
import { VendopunktoApiService } from '../../vendopunkto-api.service';
import { Action } from '@ngrx/store';
import { Router } from '@angular/router';

@Injectable()
export class InvoiceEffects implements OnInitEffects {


  loadInvoices$ = createEffect(() => {
    return this.actions$.pipe(

      ofType(InvoiceActions.loadInvoices),
      concatMap(() =>
        /** An EMPTY observable only emits completion. Replace with your own observable API request */
        this.api.searchInvoices({}).pipe(
          map(data => InvoiceActions.loadInvoicesSuccess({ invoices: data })),
          catchError(error => of(InvoiceActions.loadInvoicesFailure({ error }))))
      )
    );
  });

  loadCurrencies$ = createEffect(() => {
    return this.actions$.pipe(

      ofType(InvoiceActions.loadCurrencies),
      concatMap(() =>
        /** An EMPTY observable only emits completion. Replace with your own observable API request */
        this.api.getCurrencies().pipe(
          map(data => InvoiceActions.loadCurrenciesSuccess({ currencies: data })),
          catchError(error => of(InvoiceActions.loadCurrenciesFailure({ error }))))
      )
    );
  });

  createInvoice$ = createEffect(() => {
    return this.actions$.pipe(

      ofType(InvoiceActions.startCreateInvoice),
      concatMap((action) =>
        /** An EMPTY observable only emits completion. Replace with your own observable API request */
        this.api.createInvoice(action.invoice).pipe(
          map(data => InvoiceActions.createInvoiceSuccess({ invoice: data })),
          tap(() => this.router.navigate(['/invoices'])),
          catchError(error => of(InvoiceActions.createInvoiceFailure({ error }))))
      )
    );
  });


  ngrxOnInitEffects(): Action {
    return InvoiceActions.loadCurrencies();
  }

  constructor(
    private readonly actions$: Actions,
    private readonly api: VendopunktoApiService,
    private readonly router: Router) { }

}
