import { Component, OnInit, OnDestroy } from '@angular/core';
import { ConfigFacade } from '../+state/config.facade';
import { FormBuilder, Validators } from '@angular/forms';
import { Subject } from 'rxjs';
import { takeUntil, map, distinctUntilChanged } from 'rxjs/operators';
import { PluginsFacade } from '../../plugins/+state/plugins.facade';
import { CurrenciesFacade } from '../../currencies/+state/currencies.facade';
import * as ConfigActions from './../+state/config.actions';
import * as CurrencyActions from './../../currencies/+state/currencies.actions';
import { UpdateConfigParams } from 'shared';
import { ShellOperationsFacade } from '../../shell-operations/+state/shell-operations.facade';
import { ConfigUpdateShellOperation } from '../shell-operations';
import { createOperationInstance } from '../../shell-operations/model';
import { TitleFromCamelPipe } from '../title-from-camel.pipe';

@Component({
  selector: 'adm-config-edit-container',
  templateUrl: './config-edit-container.component.html',
  styleUrls: ['./config-edit-container.component.scss']
})
export class ConfigEditContainerComponent implements OnInit, OnDestroy {

  public readonly config$ = this.configFacade.current$;

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

  private readonly destroyedSubject = new Subject();
  private readonly titlePipe = new TitleFromCamelPipe();

  constructor(
    private readonly fb: FormBuilder,
    private readonly configFacade: ConfigFacade,
    private readonly currenciesFacade: CurrenciesFacade,
    private readonly pluginsFacade: PluginsFacade,
    private readonly shellOperations: ShellOperationsFacade,
  ) {
  }

  ngOnInit() {
    this.currenciesFacade.dispatch(CurrencyActions.loadSupportedPricingCurrencies());

    let isFirstTime = true;

    this.configFacade.current$
      .pipe(
        takeUntil(this.destroyedSubject)
      )
      .subscribe(config => {
        if (!config) {
          return;
        }

        // build a form using the config items
        for (const key in config) {
          if (config.hasOwnProperty(key)) {
            const v = config[key];
            let control = this.form.get(key);
            if (!control) {
              control = this.fb.control(v, Validators.required);
              this.form.addControl(key, control);
            }

            control.setValue(v, { emitEvent: false });

            if (isFirstTime && [
              'pricing_currencies',
              'currency_metadata_plugin',
              'exchange_rates_plugin',
              'default_pricing_currency',
            ].indexOf(key) >= 0) {
              control.valueChanges
                .pipe(
                  distinctUntilChanged(),
                  takeUntil(this.destroyedSubject)
                )
                .subscribe(newValue => {
                  this.startUpdate({ [key]: newValue });
                });
            }
          }
        }

        if (isFirstTime) {
          isFirstTime = false;
        }
      });
  }

  private startUpdate(config: UpdateConfigParams) {
    const keys = Object.keys(config)
      .map(k => this.titlePipe.transform(k))
      .join(', ');

    const op = createOperationInstance(ConfigUpdateShellOperation);
    op.title = `Configuration Update: ${keys}`;
    op.description = 'Changed ' + Object.keys(config)
      .map(k => `${this.titlePipe.transform(k)} to "${config[k]}"`);

    this.shellOperations.dispatchOperation(
      ConfigActions.updateConfigStart({ config }),
      op, true);
  }

  ngOnDestroy() {
    this.destroyedSubject.next();
    this.destroyedSubject.complete();
  }

}
