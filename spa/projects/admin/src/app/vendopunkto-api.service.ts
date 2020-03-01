import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {
  InvoiceDTO,
  InvoiceCreationParams,
  SupportedCurrency,
  GetCurrencyExchangeParams,
  GetCurrencyExchangeResult,
  GetConfigResult,
  GetPluginResult,
  UpdateConfigParams
} from 'shared';

const apiAddress = `/api/v1`;

@Injectable({
  providedIn: 'root'
})
export class VendopunktoApiService {

  constructor(private http: HttpClient) { }

  public getInvoice(invoiceId: string) {
    return this.http.get<InvoiceDTO>(`${apiAddress}/invoices/${invoiceId}`);
  }

  public searchInvoices(filter: {}) {
    return this.http.post<InvoiceDTO[]>(`${apiAddress}/invoices/search`, filter);
  }

  public getPricingCurrencies() {
    return this.http.get<SupportedCurrency[]>(`${apiAddress}/currencies/pricing`);
  }

  public searchSupportedPricingCurrencies(term?: string) {
    return this.http.post<SupportedCurrency[]>(`${apiAddress}/currencies/pricing/supported`, { term });
  }

  public getPaymentCurrencies() {
    return this.http.get<SupportedCurrency[]>(`${apiAddress}/currencies/payment-methods`);
  }

  public getCurrencyExchange(params: GetCurrencyExchangeParams) {
    return this.http.post<GetCurrencyExchangeResult>(`${apiAddress}/currencies/rates/convert`, params);
  }

  public createInvoice(params: InvoiceCreationParams) {
    return this.http.post<InvoiceDTO>(`${apiAddress}/invoices`, params);
  }

  public getConfig() {
    return this.http.get<GetConfigResult>(`${apiAddress}/config`);
  }

  public getPlugins() {
    return this.http.get<GetPluginResult[]>(`${apiAddress}/plugins`);
  }

  public updateConfig(params: UpdateConfigParams) {
    return this.http.post<GetConfigResult>(`${apiAddress}/config`, params);
  }
}
