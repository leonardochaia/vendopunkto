import { Injectable } from '@angular/core';
import { CurrenciesState } from './currencies.reducer';
import { Store, Action } from '@ngrx/store';
import {
  selectPricingCurrencies,
  selectLoadingPricingCurrencies,
  selectPaymentCurrencies,
  selectLoadingPaymentCurrencies,
  selectSupportedPricingCurrenciesCurrencies,
  selectLoadingSupportedPricingCurrencies,
} from './currencies.selectors';

@Injectable({
  providedIn: 'root'
})
export class CurrenciesFacade {

  public readonly pricingCurrencies$ = this.store.select(selectPricingCurrencies);
  public readonly loadingPricingCurrencies$ = this.store.select(selectLoadingPricingCurrencies);

  public readonly paymentCurrencies$ = this.store.select(selectPaymentCurrencies);
  public readonly loadingPaymentCurrencies$ = this.store.select(selectLoadingPaymentCurrencies);

  public readonly supportedPricingCurrencies$ = this.store.select(selectSupportedPricingCurrenciesCurrencies);
  public readonly loadingSupportedPricingCurrencies$ = this.store.select(selectLoadingSupportedPricingCurrencies);

  constructor(private readonly store: Store<CurrenciesState>) { }

  public dispatch(action: Action) {
    this.store.dispatch(action);
  }
}
