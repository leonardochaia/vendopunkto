import { createAction, props } from '@ngrx/store';
import { SupportedCurrency } from 'shared';

export const loadPricingCurrencies = createAction(
  '[Currencies] Load Pricing Currencies'
);

export const loadPricingCurrenciesSuccess = createAction(
  '[Currencies] Load Pricing Currencies Success',
  props<{ currencies: { [symbol: string]: SupportedCurrency } }>()
);

export const loadPricingCurrenciesFailure = createAction(
  '[Currencies] Load Pricing Currencies',
  props<{ error: string }>()
);

export const loadPaymentCurrencies = createAction(
  '[Currencies] Load Payment Currencies'
);

export const loadPaymentCurrenciesSuccess = createAction(
  '[Currencies] Load Payment Currencies Success',
  props<{ currencies: { [symbol: string]: SupportedCurrency } }>()
);

export const loadPaymentCurrenciesFailure = createAction(
  '[Currencies] Load Payment Currencies',
  props<{ error: string }>()
);

export const currenciesInit = createAction(
  '[Currencies] Initialization'
);