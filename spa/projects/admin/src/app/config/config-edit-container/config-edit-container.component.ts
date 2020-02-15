import { Component, OnInit, OnDestroy } from '@angular/core';
import { ConfigFacade } from '../+state/config.facade';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { Subject, combineLatest } from 'rxjs';
import { takeUntil, map, filter } from 'rxjs/operators';
import { PluginsFacade } from '../../plugins/+state/plugins.facade';
import { CurrenciesFacade } from '../../currencies/+state/currencies.facade';
import { loadPricingCurrencies } from '../../currencies/+state/currencies.actions';

@Component({
  selector: 'adm-config-edit-container',
  templateUrl: './config-edit-container.component.html',
  styleUrls: ['./config-edit-container.component.scss']
})
export class ConfigEditContainerComponent implements OnInit, OnDestroy {

  public readonly config$ = this.configFacade.current$;

  public readonly pricingCurrencies$ = this.currenciesFacade.pricingCurrencies$;

  public readonly paymentCurrencies$ = this.currenciesFacade.paymentCurrencies$;

  public readonly loadingPaymentCurrencies$ = this.currenciesFacade.loadingPaymentCurrencies$;
  public readonly loadingPricingCurrencies$ = this.currenciesFacade.loadingPricingCurrencies$;

  public readonly plugins$ = this.pluginsFacade.plugins$;

  public readonly exchangeRatesPlugins$ = this.pluginsFacade.plugins$
    .pipe(
      map(plugins => plugins.filter(p => p.pluginType === 'exchange-rate'))
    );

  public readonly currencyMetadataPlugins$ = this.pluginsFacade.plugins$
    .pipe(
      map(plugins => plugins.filter(p => p.pluginType === 'currency-metadata'))
    );

  public readonly form = this.fb.group({});

  public get pricingCurrenciesControl() {
    return this.form.get('pricing_currencies') as FormGroup;
  }

  public get paymentMethodsControl() {
    return this.form.get('payment_methods') as FormGroup;
  }


  private readonly destroyedSubject = new Subject();

  constructor(
    private readonly fb: FormBuilder,
    private readonly configFacade: ConfigFacade,
    private readonly currenciesFacade: CurrenciesFacade,
    private readonly pluginsFacade: PluginsFacade,
  ) {
  }

  ngOnInit() {
    combineLatest(this.configFacade.current$,
      this.pricingCurrencies$.pipe(filter(c => !!c)),
      this.paymentCurrencies$.pipe(filter(c => !!c))
    )
      .pipe(
        takeUntil(this.destroyedSubject)
      )
      .subscribe(([config, pricingCurrencies, paymentCurrencies]) => {
        if (!config) {
          return;
        }
        for (const key in config) {
          if (config.hasOwnProperty(key)) {
            const v = config[key];
            let control = this.form.get(key);
            if (!control) {
              control = this.fb.control(v, Validators.required);
              this.form.addControl(key, control);
            }

            control.setValue(v);
          }
        }

        const pricingCurrenciesGroup = this.fb.group({});

        for (const key in pricingCurrencies) {
          if (pricingCurrencies.hasOwnProperty(key)) {
            const currency = pricingCurrencies[key];
            const enabled = (config.pricing_currencies as string[])
              .indexOf(currency.symbol.toLowerCase()) >= 0;
            const control = this.fb.control(enabled);
            pricingCurrenciesGroup.addControl(currency.symbol, control);
          }
        }

        this.form.setControl('pricing_currencies', pricingCurrenciesGroup);

        const paymentMethodsGroup = this.fb.group({});
        for (const key in paymentCurrencies) {
          if (paymentCurrencies.hasOwnProperty(key)) {
            const currency = paymentCurrencies[key];
            const enabled = (config.payment_methods as string[])
              .indexOf(currency.symbol.toLowerCase()) >= 0;
            const control = this.fb.control(enabled);
            paymentMethodsGroup.addControl(currency.symbol, control);
          }
        }
        this.form.setControl('payment_methods', paymentMethodsGroup);
      });
  }

  public removePricingCurrency(key: string) {
    this.pricingCurrenciesControl.removeControl(key);
  }

  ngOnDestroy() {
    this.destroyedSubject.next();
    this.destroyedSubject.complete();
  }

}
