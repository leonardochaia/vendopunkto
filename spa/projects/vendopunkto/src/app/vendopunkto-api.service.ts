import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { InvoiceDTO } from './model';
import { webSocket } from 'rxjs/webSocket';

const apiAddress = `/api/v1`;

@Injectable({
  providedIn: 'root'
})
export class VendopunktoApiService {

  constructor(private http: HttpClient) { }

  public getInvoice(invoiceId: string) {
    return this.http.get<InvoiceDTO>(`${apiAddress}/invoices/${invoiceId}`);
  }

  public buildInvoiceSocket(invoiceId: string) {
    const proto = (window.location.protocol === 'https:' ? 'wss://' : 'ws://');
    return webSocket<InvoiceDTO>(proto + window.location.host + `${apiAddress}/invoices/${invoiceId}/ws`);
  }

  public generatePaymentMethodAddress(invoiceId: string, currency: string) {
    const url = `${apiAddress}/invoices/${invoiceId}/payment-method/address`;
    const params = {
      currency
    };
    return this.http.post<InvoiceDTO>(url, params);
  }
}
