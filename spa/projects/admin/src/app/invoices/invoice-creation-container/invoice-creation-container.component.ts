import { Component, OnInit, OnDestroy } from '@angular/core';
import { InvoiceFacade } from '../+state/invoice.facade';
import { FormBuilder, Validators } from '@angular/forms';
import { startCreateInvoice, invoiceCreationFormChanged } from '../+state/invoice.actions';
import { map, takeUntil, filter, debounceTime, distinctUntilChanged, startWith } from 'rxjs/operators';
import { InvoiceCreationParams, PaymentMethodCreationParams, SupportedCurrency } from 'shared';
import { Subject, combineLatest } from 'rxjs';
import { CurrenciesFacade } from '../../currencies/+state/currencies.facade';
import { MatSelectChange } from '@angular/material';

@Component({
  selector: 'adm-invoice-creation-container',
  templateUrl: './invoice-creation-container.component.html',
  styleUrls: ['./invoice-creation-container.component.scss']
})
export class InvoiceCreationContainerComponent implements OnInit, OnDestroy {

  public readonly totalControl = this.fb.control(null, Validators.required);
  public readonly currencyControl = this.fb.control(null, Validators.required);
  public readonly paymentMethodsArray = this.fb.array([]);

  public readonly basicInfoForm = this.fb.group({
    total: this.totalControl,
    currency: this.currencyControl,
    paymentMethods: this.paymentMethodsArray
  });

  public readonly pricingCurrencies$ = this.currenciesFacade.pricingCurrencies$;
  public readonly paymentCurrencies$ = this.currenciesFacade.paymentCurrencies$
    .pipe(
      map(arr => arr.reduce((a, b) => (a[b.symbol] = b, a), {})),
      map(o => o as { [symbol: string]: SupportedCurrency })
    );

  public readonly loadingPricingCurrencies$ = this.currenciesFacade.loadingPricingCurrencies$;
  public readonly loadingPaymentCurrencies$ = this.currenciesFacade.loadingPaymentCurrencies$;
  public readonly loadingPaymentMethods$ = this.facade.loadingPaymentMethods$
    .pipe(
      debounceTime(250),
      distinctUntilChanged(),
    );
  public readonly creating$ = this.facade.creating$;

  public readonly currentPricingCurrency$ = combineLatest(this.currencyControl.valueChanges
    .pipe(
      distinctUntilChanged(),
      map(v => v as string),
      startWith(null),
    ), this.pricingCurrencies$)
    .pipe(
      map(([symbol, pricingCurrencies]) => pricingCurrencies.find(c => c.symbol === symbol)),
    );


  private readonly destroyedSubject = new Subject();

  constructor(
    private readonly facade: InvoiceFacade,
    private readonly currenciesFacade: CurrenciesFacade,
    private readonly fb: FormBuilder) { }

  ngOnInit() {

    this.facade.creationBasicInfo$
      .pipe(
        takeUntil(this.destroyedSubject)
      )
      .subscribe(info => {
        if (!this.basicInfoForm.dirty) {
          this.currencyControl.setValue(info.currency);
          this.totalControl.setValue(info.total);
        }

        this.paymentMethodsArray.clear();
        for (const pm of info.paymentMethods) {
          const group = this.fb.group({
            currency: [pm.currency, Validators.required],
            total: [pm.total, Validators.required]
          });

          if (pm.currency === info.currency) {
            group.disable();
            this.paymentMethodsArray.insert(0, group);
          } else {
            this.paymentMethodsArray.push(group);
          }
        }
      });

    const totalChanges$ = this.totalControl.valueChanges
      .pipe(
        distinctUntilChanged(),
        map(v => v as number),
        startWith(null),
      );

    const currencyChanges$ = this.currencyControl.valueChanges
      .pipe(
        distinctUntilChanged(),
        map(v => v as string),
        startWith(null),
      );

    combineLatest(totalChanges$, currencyChanges$)
      .pipe(
        debounceTime(250),
        filter(() => this.basicInfoForm.valid),
        map(() => this.basicInfoForm.value as InvoiceCreationParams),
        takeUntil(this.destroyedSubject)
      ).subscribe(form => {
        this.facade.dispatch(invoiceCreationFormChanged({ form }));
      });
  }

  ngOnDestroy() {
    this.destroyedSubject.next();
    this.destroyedSubject.complete();
  }

  public createInvoice() {

    const params: InvoiceCreationParams = this.basicInfoForm.getRawValue();

    params.paymentMethods = this.paymentMethodsArray.getRawValue()
      .map(c => ({ currency: c.currency, total: c.total } as PaymentMethodCreationParams));

    this.facade.dispatch(startCreateInvoice({
      invoice: params
    }));
  }
}
