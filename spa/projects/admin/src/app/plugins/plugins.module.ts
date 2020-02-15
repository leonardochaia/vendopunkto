import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';

import * as fromConfig from './+state/plugins.reducer';
import { PluginEffects } from './+state/plugins.effects';


@NgModule({
  declarations: [],
  imports: [
    CommonModule,

    StoreModule.forFeature(fromConfig.PluginsFeatureKey, fromConfig.reducer),
    EffectsModule.forFeature([PluginEffects]),
  ]
})
export class PluginsModule { }
