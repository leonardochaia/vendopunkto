import { createAction, props } from '@ngrx/store';
import { InvoiceDTO, SupportedCurrency, InvoiceCreationParams } from 'shared';

export const loadInvoices = createAction(
  '[Invoice] Load Invoices'
);

export const loadInvoicesSuccess = createAction(
  '[Invoice] Load Invoices Success',
  props<{ invoices: InvoiceDTO[] }>()
);

export const loadInvoicesFailure = createAction(
  '[Invoice] Load Invoices Failure',
  props<{ error: Error }>()
);

export const startCreateInvoice = createAction(
  '[Invoice] Start Create Invoice',
  props<{ invoice: InvoiceCreationParams }>()
);

export const createInvoiceSuccess = createAction(
  '[Invoice] Create Invoice Success',
  props<{ invoice: InvoiceDTO }>()
);

export const createInvoiceFailure = createAction(
  '[Invoice] Create Invoice Failure',
  props<{ error: Error }>()
);


export const loadCurrencies = createAction(
  '[Invoice] Load Currencies'
);

export const loadCurrenciesSuccess = createAction(
  '[Invoice] Load Currencies Success',
  props<{ currencies: SupportedCurrency[] }>()
);

export const loadCurrenciesFailure = createAction(
  '[Invoice] Load Currencies Failure',
  props<{ error: Error }>()
);
