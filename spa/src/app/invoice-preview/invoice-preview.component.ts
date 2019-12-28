import { Component, OnInit, Input, OnDestroy } from '@angular/core';
import { takeUntil, map, withLatestFrom, shareReplay, startWith } from 'rxjs/operators';
import { InvoiceDTO } from '../model';
import { VendopunktoApiService } from '../vendopunkto-api.service';
import { Subject, ReplaySubject, combineLatest } from 'rxjs';

@Component({
  selector: 'app-invoice-preview',
  templateUrl: './invoice-preview.component.html',
  styleUrls: ['./invoice-preview.component.scss']
})
export class InvoicePreviewComponent implements OnInit, OnDestroy {

  protected readonly invoiceSubject = new ReplaySubject<InvoiceDTO>(1);
  protected readonly currencySubject = new Subject<string>();
  protected readonly destroyedSubject = new Subject();

  public readonly webSocketSupported = 'WebSocket' in window;

  public readonly invoice$ = this.invoiceSubject.asObservable();

  @Input()
  public invoiceId: string;

  public readonly paymentMethod$ = combineLatest(this.invoice$, this.currencySubject.pipe(startWith(null)))
    .pipe(
      map(([invoice, currency]) => invoice.paymentMethods
        .filter(pm => pm.currency === (currency || invoice.currency))[0])
    );

  constructor(private readonly api: VendopunktoApiService) { }

  public ngOnInit() {
    this.initializeInvoice();
  }

  public updateInvoice() {
    this.api.getInvoice(this.invoiceId)
      .subscribe(freshInvoice => {
        this.invoiceSubject.next(freshInvoice);
      });
  }

  public changePaymentMethod(currency: string) {
    this.currencySubject.next(currency);
  }

  public ngOnDestroy() {
    this.destroyedSubject.next();
    this.destroyedSubject.complete();
  }

  protected initializeInvoice() {
    if (this.webSocketSupported) {
      this.api.buildInvoiceSocket(this.invoiceId)
        .pipe(
          shareReplay(),
          takeUntil(this.destroyedSubject)
        ).subscribe(invoice => {
          console.log(invoice);
          // this wraps the subject
          // to prevent exposing the actual socket to the template
          this.invoiceSubject.next(invoice);
        }, err => {
          console.error(err);
        });
    } else {
      console.warn('WebSocket not supported. Fallback to polling.');
      this.updateInvoice();
    }

    this.currencySubject.pipe(
      takeUntil(this.destroyedSubject),
      withLatestFrom(this.invoice$)
    )
      .subscribe(([currency, invoice]) => {
        const method = invoice.paymentMethods.filter(p => p.currency === currency)[0];
        if (!method.address) {
          this.api.generatePaymentMethodAddress(this.invoiceId, currency)
            .subscribe();
        }
      });
  }

}
