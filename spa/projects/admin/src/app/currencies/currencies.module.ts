import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';

import * as fromCurrencies from './+state/currencies.reducer';
import { CurrenciesEffects } from './+state/currencies.effects';


@NgModule({
  declarations: [],
  imports: [
    CommonModule,

    StoreModule.forFeature(fromCurrencies.CurrenciesFeatureKey, fromCurrencies.reducer),
    EffectsModule.forFeature([CurrenciesEffects]),
  ]
})
export class CurrenciesModule { }
