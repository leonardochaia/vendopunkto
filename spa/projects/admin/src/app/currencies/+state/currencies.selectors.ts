import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromCurrencies from './currencies.reducer';

export const selectCurrenciesState = createFeatureSelector<fromCurrencies.CurrenciesState>(
  fromCurrencies.CurrenciesFeatureKey
);

export const selectPricingCurrencies = createSelector(selectCurrenciesState,
  s => s.pricingCurrencies);

export const selectLoadingPricingCurrencies = createSelector(selectCurrenciesState,
  s => s.loadingPricingCurrencies);

export const selectPaymentCurrencies = createSelector(selectCurrenciesState,
  s => s.paymentCurrencies);

export const selectLoadingPaymentCurrencies = createSelector(selectCurrenciesState,
  s => s.loadingPaymentCurrencies);
