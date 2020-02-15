import { Action, createReducer, on } from '@ngrx/store';
import * as currenciesActions from './currencies.actions';
import { SupportedCurrency } from 'shared';

export const CurrenciesFeatureKey = 'vpCurrencies';

export interface CurrenciesState {
  pricingCurrencies: { [symbol: string]: SupportedCurrency };
  loadingPricingCurrencies: boolean;

  paymentCurrencies: { [symbol: string]: SupportedCurrency };
  loadingPaymentCurrencies: boolean;

  error: string;
}

export const initialState: CurrenciesState = {
  pricingCurrencies: undefined,
  loadingPricingCurrencies: false,
  paymentCurrencies: undefined,
  loadingPaymentCurrencies: false,
  error: null
};

const CurrenciesReducer = createReducer(
  initialState,

  on(currenciesActions.loadPricingCurrencies, (state) => ({
    ...state,
    loadingPricingCurrencies: true
  })),

  on(currenciesActions.loadPricingCurrenciesSuccess, (state, action) => ({
    ...state,
    loadingPricingCurrencies: false,
    pricingCurrencies: action.currencies
  })),

  on(currenciesActions.loadPricingCurrenciesFailure, (state, action) => ({
    ...state,
    loadingPricingCurrencies: false,
    error: action.error
  })),

  on(currenciesActions.loadPaymentCurrencies, (state) => ({
    ...state,
    loadingPaymentCurrencies: true
  })),

  on(currenciesActions.loadPaymentCurrenciesSuccess, (state, action) => ({
    ...state,
    loadingPaymentCurrencies: false,
    paymentCurrencies: action.currencies
  })),

  on(currenciesActions.loadPaymentCurrenciesFailure, (state, action) => ({
    ...state,
    loadingPaymentCurrencies: false,
    error: action.error
  })),
);

export function reducer(state: CurrenciesState | undefined, action: Action) {
  return CurrenciesReducer(state, action);
}
