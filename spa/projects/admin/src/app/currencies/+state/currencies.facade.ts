import { Injectable } from '@angular/core';
import { CurrenciesState } from './currencies.reducer';
import { Store, Action } from '@ngrx/store';
import {
  selectPricingCurrencies, selectLoadingPricingCurrencies, selectPaymentCurrencies, selectLoadingPaymentCurrencies,
} from './currencies.selectors';

@Injectable({
  providedIn: 'root'
})
export class CurrenciesFacade {

  public readonly pricingCurrencies$ = this.store.select(selectPricingCurrencies);
  public readonly loadingPricingCurrencies$ = this.store.select(selectLoadingPricingCurrencies);

  public readonly paymentCurrencies$ = this.store.select(selectPaymentCurrencies);
  public readonly loadingPaymentCurrencies$ = this.store.select(selectLoadingPaymentCurrencies);
  // public readonly pricingCurrencies$ = this.store.select(selectPricingCurrencies);

  constructor(private readonly store: Store<CurrenciesState>) { }

  public dispatch(action: Action) {
    this.store.dispatch(action);
  }
}
