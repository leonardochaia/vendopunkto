import { Component, Input } from '@angular/core';
import { CurrenciesFacade } from '../../currencies/+state/currencies.facade';
import { filter, map } from 'rxjs/operators';
import { ShellSelectorItem } from '../../shell-selector/shell-selector-dialog/shell-selector-dialog.component';
import { FormGroup } from '@angular/forms';

@Component({
  selector: 'adm-config-invoice-edit',
  templateUrl: './config-invoice-edit.component.html',
  styleUrls: ['./config-invoice-edit.component.scss']
})
export class ConfigInvoiceEditComponent {

  @Input()
  public form: FormGroup;

  public readonly pricingCurrencies$ = this.currenciesFacade.pricingCurrencies$;

  public readonly supportedPricingCurrencies$ = this.currenciesFacade.supportedPricingCurrencies$.pipe(
    filter(c => !!c),
    // map to the shell-selector items
    map(c => Object.keys(c).map(k => ({
      id: c[k].symbol,
      name: `${c[k].name} (${c[k].symbol.toUpperCase()})`,
      initials: c[k].symbol,
      imageURL: c[k].logoImageUrl,
    } as ShellSelectorItem)))
  );

  constructor(
    private readonly currenciesFacade: CurrenciesFacade,
  ) { }

}
