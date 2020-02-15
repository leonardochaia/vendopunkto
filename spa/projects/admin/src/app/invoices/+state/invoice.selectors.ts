import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromInvoice from './invoice.reducer';
import { SupportedCurrency } from 'shared';

export const selectInvoiceState = createFeatureSelector<fromInvoice.InvoiceState>(
  fromInvoice.InvoiceFeatureKey
);

export const selectLoadingInvoices = createSelector(selectInvoiceState, s => s.loadingInvoices);

export const selectInvoices = createSelector(selectInvoiceState, s => s.invoices);

export const selectLoadingPaymentMethods = createSelector(selectInvoiceState, s => s.creation.loadingPaymentMethods);

export const selectBasicCreationInfo = createSelector(selectInvoiceState, s => s.creation.basic);

export const selectCreatingInvoice = createSelector(selectInvoiceState, s => s.creation.creating);
