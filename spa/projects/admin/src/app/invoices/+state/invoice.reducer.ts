import { Action, createReducer, on } from '@ngrx/store';
import * as InvoiceActions from './invoice.actions';
import { InvoiceDTO, SupportedCurrency } from 'shared';

export const InvoiceFeatureKey = 'vpInvoice';

export interface InvoiceState {
  invoices: InvoiceDTO[];
  currencies: SupportedCurrency[];
  loading: boolean;
}

export const initialState: InvoiceState = {
  invoices: [],
  currencies: [],
  loading: true
};

const InvoiceReducer = createReducer(
  initialState,

  on(InvoiceActions.loadInvoices, state => ({
    ...state,
    loading: true
  })),
  on(InvoiceActions.loadInvoicesSuccess, (state, action) => ({
    ...state,
    loading: false,
    invoices: action.invoices,
  })),
  on(InvoiceActions.loadInvoicesFailure, (state, action) => ({
    ...state,
    loading: false
  })),

  on(InvoiceActions.loadCurrenciesSuccess, (state, action) => ({
    ...state,
    currencies: action.currencies,
  })),
);

export function reducer(state: InvoiceState | undefined, action: Action) {
  return InvoiceReducer(state, action);
}
