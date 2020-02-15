import { Component, OnInit, Input, TemplateRef } from '@angular/core';
import { SupportedCurrency } from 'shared/shared';

@Component({
  selector: 'adm-currency-card',
  templateUrl: './currency-card.component.html',
  styleUrls: ['./currency-card.component.scss']
})
export class CurrencyCardComponent implements OnInit {

  @Input()
  public currency: SupportedCurrency;

  @Input()
  public actionTemplate: TemplateRef<SupportedCurrency>;

  constructor() { }

  ngOnInit() {
  }

}
