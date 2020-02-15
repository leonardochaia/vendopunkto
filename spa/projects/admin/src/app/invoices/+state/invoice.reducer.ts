import { Action, createReducer, on } from '@ngrx/store';
import * as InvoiceActions from './invoice.actions';
import * as CurrenciesActions from '../../currencies/+state/currencies.actions';
import * as ConfigActions from '../../config/+state/config.actions';
import {
  InvoiceDTO,
  InvoiceCreationParams,
  PaymentMethodCreationParams
} from 'shared';

export const InvoiceFeatureKey = 'vpInvoice';

export interface InvoiceState {
  error: string;

  invoices: InvoiceDTO[];
  loadingInvoices: boolean;

  creation: {
    basic: InvoiceCreationParams;
    loadingPaymentMethods: boolean;
    creating: boolean;
  };
}

export const initialState: InvoiceState = {
  invoices: [],
  loadingInvoices: false,

  creation: {
    basic: {
      currency: null,
      paymentMethods: [],
      total: null,
    },
    creating: false,
    loadingPaymentMethods: false
  },
  error: null,
};

const InvoiceReducer = createReducer(
  initialState,

  on(InvoiceActions.loadInvoices, state => ({
    ...state,
    loadingInvoices: true
  })),
  on(InvoiceActions.loadInvoicesSuccess, (state, action) => ({
    ...state,
    loadingInvoices: false,
    invoices: action.invoices,
  })),
  on(InvoiceActions.loadInvoicesFailure, (state, action) => ({
    ...state,
    loadingInvoices: false
  })),

  on(ConfigActions.loadConfigSuccess, (state, action) => ({
    ...state,
    creation: {
      ...state.creation,
      basic: {
        ...state.creation.basic,
        currency: action.config.default_pricing_currency
      }
    },
    loadingCurrencies: false
  })),

  on(InvoiceActions.invoiceCreationFormChanged, (state, action) => ({
    ...state,
    creation: {
      ...state.creation,
      loadingPaymentMethods: true,
      basic: {
        ...state.creation.basic,
        ...action.form,
      },
    }
  })),

  on(InvoiceActions.getPaymentMethodExchangeRateSuccess, (state, action) => ({
    ...state,
    creation: {
      ...state.creation,
      loadingPaymentMethods: false,
      basic: {
        ...state.creation.basic,

        paymentMethods: Object.keys(action.result).map(k => ({
          currency: k,
          total: action.result[k] as any,
        } as PaymentMethodCreationParams))
      }
    }
  })),

  on(InvoiceActions.getPaymentMethodExchangeRateFailure, (state, action) => ({
    ...state,
    creation: {
      ...state.creation,
      loadingPaymentMethods: false,
    },
    error: action.error
  })),

  on(InvoiceActions.startCreateInvoice, (state, action) => ({
    ...state,
    creation: {
      ...state.creation,
      creating: true
    },
  })),

  on(InvoiceActions.createInvoiceSuccess, (state, action) => ({
    ...state,
    creation: {
      ...state.creation,
      loadingPaymentMethods: false,
      basic: {
        ...state.creation.basic,
        paymentMethods: [],
        total: null,
      },
      creating: false
    },
  })),

  on(InvoiceActions.createInvoiceFailure, (state, action) => ({
    ...state,
    creation: {
      ...state.creation,
      creating: false
    },
    error: action.error,
  })),
);

export function reducer(state: InvoiceState | undefined, action: Action) {
  return InvoiceReducer(state, action);
}
