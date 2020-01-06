import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { InvoiceDTO, InvoiceCreationParams, SupportedCurrency } from 'shared';

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

  public getCurrencies() {
    return this.http.get<SupportedCurrency[]>(`${apiAddress}/currencies`);
  }

  public createInvoice(params: InvoiceCreationParams) {
    return this.http.post<InvoiceDTO>(`${apiAddress}/invoices`, params);
  }
}
