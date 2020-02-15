import { createAction, props } from '@ngrx/store';
import { InvoiceDTO, SupportedCurrency, InvoiceCreationParams, GetCurrencyExchangeResult } from 'shared';

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
  props<{ error: string }>()
);

export const invoiceCreationFormChanged = createAction(
  '[Invoice] Creation Form Changed',
  props<{ form: InvoiceCreationParams }>()
);


export const getPaymentMethodExchangeRateSuccess = createAction(
  '[Invoice] Get PaymentMethod Exchange Rate Success',
  props<{ result: GetCurrencyExchangeResult }>()
);

export const getPaymentMethodExchangeRateFailure = createAction(
  '[Invoice] Get PaymentMethod Exchange Rate Failure',
  props<{ error: string }>()
);

