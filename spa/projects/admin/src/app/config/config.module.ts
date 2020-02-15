import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { ConfigEditContainerComponent } from './config-edit-container/config-edit-container.component';
import {
  MatTabsModule,
  MatFormFieldModule,
  MatInputModule,
  MatButtonModule,
  MatIconModule,
  MatSelectModule,
  MatCardModule,
  MatCheckboxModule
} from '@angular/material';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ReactiveFormsModule } from '@angular/forms';

import * as fromConfig from './+state/config.reducer';
import { ConfigEffects } from './+state/config.effects';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { CurrencyCardComponent } from './currency-card/currency-card.component';
import { ConfigSelectorCardComponent } from './config-selector-card/config-selector-card.component';

const routes: Routes = [
  {
    path: 'config',
    component: ConfigEditContainerComponent,
  },
]

@NgModule({
  declarations: [
    ConfigEditContainerComponent,
    CurrencyCardComponent,
    ConfigSelectorCardComponent
  ],
  imports: [
    CommonModule,
    RouterModule.forChild(routes),
    ReactiveFormsModule,

    FlexLayoutModule,
    MatTabsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatSelectModule,
    MatCardModule,
    MatCheckboxModule,

    StoreModule.forFeature(fromConfig.ConfigFeatureKey, fromConfig.reducer),
    EffectsModule.forFeature([ConfigEffects]),

  ]
})
export class ConfigModule { }
