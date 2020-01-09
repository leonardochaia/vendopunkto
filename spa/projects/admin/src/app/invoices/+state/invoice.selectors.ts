import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromInvoice from './invoice.reducer';
import { SupportedCurrency } from 'shared';

export const selectInvoiceState = createFeatureSelector<fromInvoice.InvoiceState>(
  fromInvoice.InvoiceFeatureKey
);

export const selectLoadingInvoices = createSelector(selectInvoiceState, s => s.loadingInvoices);

export const selectInvoices = createSelector(selectInvoiceState, s => s.invoices);


export const selectCurrencies = createSelector(selectInvoiceState, s => s.currencies);

export const selectPaymentCurrencies = createSelector(selectInvoiceState, s => {
  const output: SupportedCurrency[] = [];
  for (const k in s.currencies) {
    if (s.currencies.hasOwnProperty(k)) {
      const value = s.currencies[k];
      if (value.supportsPayments) {
        output.push(value);
      }
    }
  }
  return output;
});

export const selectLoadingCurrencies = createSelector(selectInvoiceState, s => s.loadingCurrencies);


export const selectLoadingPaymentMethods = createSelector(selectInvoiceState, s => s.creation.loadingPaymentMethods);

export const selectBasicCreationInfo = createSelector(selectInvoiceState, s => s.creation.basic);

export const selectCreatingInvoice = createSelector(selectInvoiceState, s => s.creation.creating);
