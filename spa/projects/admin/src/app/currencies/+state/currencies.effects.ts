import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, concatMap, switchMap } from 'rxjs/operators';
import { of } from 'rxjs';

import * as currenciesActions from './currencies.actions';
import { Action } from '@ngrx/store';
import { VendopunktoApiService } from '../../vendopunkto-api.service';

@Injectable()
export class CurrenciesEffects {

  onLoadPricingCurrencies$ = createEffect(() => this.actions$
    .pipe(
      ofType(currenciesActions.loadPricingCurrencies),

      concatMap(() => this.api.getPricingCurrencies()
        .pipe(
          map(data => currenciesActions.loadPricingCurrenciesSuccess({
            currencies: data
          })),
          catchError(e => of(currenciesActions.loadPricingCurrenciesFailure({ error: e.message })))
        ))
    ));

  onLoadPaymentCurrencies$ = createEffect(() => this.actions$
    .pipe(
      ofType(currenciesActions.loadPaymentCurrencies),

      concatMap(() => this.api.getPaymentCurrencies()
        .pipe(
          map(data => currenciesActions.loadPaymentCurrenciesSuccess({
            currencies: data
          })),
          catchError(e => of(currenciesActions.loadPaymentCurrenciesFailure({ error: e.message })))
        ))
    ));

  onLoadSupportedPricingCurrencies$ = createEffect(() => this.actions$
    .pipe(
      ofType(currenciesActions.loadSupportedPricingCurrencies),

      concatMap(() => this.api.searchSupportedPricingCurrencies()
        .pipe(
          map(data => currenciesActions.loadSupportedPricingCurrenciesSuccess({
            currencies: data,
          })),
          catchError(e => of(currenciesActions.loadSupportedPricingCurrenciesFailure({ error: e.message })))
        ))
    ));

  onInit$ = createEffect(() => this.actions$
    .pipe(
      ofType(currenciesActions.currenciesInit),

      switchMap(() => [
        currenciesActions.loadPricingCurrencies(),
        currenciesActions.loadPaymentCurrencies(),
      ])
    ));

  ngrxOnInitEffects(): Action {
    return currenciesActions.currenciesInit();
  }

  constructor(
    private readonly actions$: Actions,
    private readonly api: VendopunktoApiService) { }

}
