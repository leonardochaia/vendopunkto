import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';

import * as fromShellOperations from './+state/shell-operations.reducer';
import { ShellOperationsEffects } from './+state/shell-operations.effects';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,

    StoreModule.forFeature(fromShellOperations.ShellOperationsFeatureKey,
      fromShellOperations.reducer),
    EffectsModule.forFeature([ShellOperationsEffects]),

  ]
})
export class ShellOperationsModule { }
