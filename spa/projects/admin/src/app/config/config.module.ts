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
  MatCheckboxModule,
  MatProgressBarModule,
} from '@angular/material';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ReactiveFormsModule } from '@angular/forms';

import * as fromConfig from './+state/config.reducer';
import { ConfigEffects } from './+state/config.effects';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { CurrencyCardComponent } from './currency-card/currency-card.component';
import { ConfigSelectorCardComponent } from './config-selector-card/config-selector-card.component';
import { ShellSelectorModule } from '../shell-selector/shell-selector.module';
import { ShellAvatarModule } from '../shell-avatar/shell-avatar.module';
import { ConfigInvoiceEditComponent } from './config-invoice-edit/config-invoice-edit.component';
import { TitleFromCamelPipe } from './title-from-camel.pipe';

const routes: Routes = [
  {
    path: 'config',
    component: ConfigEditContainerComponent,
  },
];

@NgModule({
  declarations: [
    ConfigEditContainerComponent,
    CurrencyCardComponent,
    ConfigSelectorCardComponent,
    ConfigInvoiceEditComponent,
    TitleFromCamelPipe
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
    MatProgressBarModule,

    StoreModule.forFeature(fromConfig.ConfigFeatureKey, fromConfig.reducer),
    EffectsModule.forFeature([ConfigEffects]),

    ShellSelectorModule,
    ShellAvatarModule,

  ]
})
export class ConfigModule { }
