import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromInvoice from './invoice.reducer';

export const selectInvoiceState = createFeatureSelector<fromInvoice.InvoiceState>(
  fromInvoice.InvoiceFeatureKey
);

export const selectInvoices = createSelector(selectInvoiceState, s => s.invoices);

export const selectCurrencies = createSelector(selectInvoiceState, s => s.currencies);
